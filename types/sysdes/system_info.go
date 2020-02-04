package sysdes

import (
	"bufio"
	"bytes"
	"github.com/geekymedic/neon-cli/util"
	"github.com/geekymedic/neon/logger"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon/errors"
)

const (
	SystemNameSubffix = "-system"
	ProtocolName      = "protocol"
	GoHostPreffix     = "git.gmtshenzhen.com/yaodao"
)

type SystemDes struct {
	Name       string `yaml:"name"` // "foo-system"
	ShortName  string
	Author     string `yaml:"author"`
	GoModel    string `yaml:"gomod"`
	CreateTime string `yaml:"create_time"`
	UpdateTime string `yaml:"update_time"`
	Bffs       *Bffs  `yaml:"bff,omitempty"`
	//Services   Services `yaml:"services,omitempty"`
	DirNode types.DirNode `json:",omitempty" yaml:"-"`
}

func NewSystemDes(dirNode interface{}, bffName ...string) (*SystemDes, error) {
	var sys = &SystemDes{
		Author:     "GeekyMedic",
		CreateTime: time.Now().String(),
		UpdateTime: time.Now().String(),
	}
	switch typ := dirNode.(type) {
	case string:
		dir := types.NewBaseDir(typ)
		baseDir := regexp.MustCompile("(.*system)*").FindString(dir.Abs())
		if baseDir == "" {
			return nil, errors.Format("System direcoty is invalid: %s", dir.Abs())
		}
		dir = types.NewBaseDir(baseDir)
		sys.DirNode = dir
		dirs := dir.Split()
		sys.Name = dirs[len(dirs)-1]
	case types.DirNode:
		dir := typ
		baseDir := regexp.MustCompile("(.*system)*").FindString(dir.Abs())
		if baseDir == "" {
			return nil, errors.Format("System direcoty is invalid: %s", dir.Abs())
		}
		dir = types.NewBaseDir(baseDir)
		sys.DirNode = dir
		dirs := dir.Split()
		sys.Name = dirs[len(dirs)-1]
	default:
		types.PanicSanityf("Unsupport type:%T", dirNode)
	}

	idx := strings.LastIndex(sys.Name, "-")
	sys.ShortName = sys.Name[0:idx]
	bffs, err := NewBffs(sys, bffName...)
	if err != nil {
		return nil, err
	}
	if bffs == nil {
		return nil, errors.NewStackError("not found any bff info")
	}
	sys.Bffs = bffs
	gomod := util.ConvertBreakLinePath(sys.DirNode.Abs() + "/"+ "go.mod")
	gomodFp := types.NewBaseFile(gomod)
	types.AssertNil(gomodFp.Create(os.O_RDONLY, os.ModePerm))
	defer gomodFp.Close()
	buf, err := gomodFp.ReadAll()
	if err != nil {
		return nil, err
	}
	line, _ , err := bufio.NewReader(bytes.NewBuffer(buf)).ReadLine()
	if err != nil {
		return  nil, err
	}
	lines := strings.Split(string(line), " ")
	sys.GoModel = lines[len(lines)-1]
	return sys, nil
}

type Bffs struct {
	DirNode  types.DirNode `json:",omitempty" yaml:"-"`
	BffItems []*BffItem    `json:",omitempty" yaml:"names"`
	Sys      *SystemDes    `json:",omitempty" yaml:"-"`
}

func NewBffs(sys *SystemDes, bffName ...string) (*Bffs, error) {
	var (
		bffDir = sys.DirNode.Append("bff").(types.DirNode)
		items  []*BffItem
		err    error
	)
	err = bffDir.Walk(func(path string, info os.FileInfo, err error) error {
		if info == nil || !info.IsDir() || path == bffDir.Abs() || err != nil {
			return nil
		}
		if len(bffName) > 0 {
			if strings.Contains(path, bffName[0]) {
				logger.With("path", path).Info("Found a bff")
			} else {
				return nil
			}
		}
		item, err := NewBffItem(types.NewBaseDir(path), sys)
		if err != nil {
			logger.Warnf("Fail to create bff: %s", path)
		} else if item != nil && len(item.Impls) > 0 {
			items = append(items, item)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	// if len(items) == 0 {
	// 	return nil, errors.NewStackError("Not found any bff info")
	// }
	return &Bffs{DirNode: bffDir, BffItems: items, Sys: sys}, nil
}

func (bff Bffs) BffIter(iterFn func(item *BffItem) bool) {
	for _, item := range bff.BffItems {
		if !iterFn(item) {
			return
		}
	}
}

func (bff Bffs) ImplIter(iterFn func(*BffItem, *BffImpl) bool) {
	for _, item := range bff.BffItems {
		for _, impl := range item.Impls {
			if !iterFn(item, impl) {
				return
			}
		}
	}
}

func (bff Bffs) MatchBff(fileNode types.FileNode) *BffItem {
	for _, bff := range bff.BffItems {
		if strings.HasPrefix(fileNode.Abs(), bff.DirNode.Abs()) {
			return bff
		}
	}
	return nil
}

// eg: MatchBffImplByPath("/xxx/user_system/bff/admin/impls/login.go"
func (bff Bffs) MatchBffImplByPath(fileNode types.FileNode) *BffImpl {
	bffItem := bff.MatchBff(fileNode)
	if bffItem == nil {
		return nil
	}
	return bffItem.MatchImpl(fileNode)
}

func (bff Bffs) MatchBffByName(bffName string) *BffItem {
	for _, bff := range bff.BffItems {
		if bff.DirNode.Name() == bffName {
			return bff
		}
	}
	return nil
}

func (bff Bffs) MatchBffImplByName(bffName, implName string) *BffImpl {
	bffItem := bff.MatchBffByName(bffName)
	if bffItem == nil {
		return nil
	}
	return bffItem.MatchImplByName(implName)
}

type BffItem struct {
	DirNode types.DirNode `json:",omitempty" yaml:"-"`
	Impls   []*BffImpl    `json:",omitempty" yaml:"interfaces"`
	Sys     *SystemDes    `json:",omitempty" yaml:"-"`
}

func NewBffItem(dirNode types.DirNode, sys *SystemDes) (*BffItem, error) {
	var bff = new(BffItem)
	bff.DirNode = dirNode
	implAbs := types.NewBaseDir(dirNode.Append("impls").(types.DirNode).Abs())
	err := filepath.Walk(implAbs.Abs(), func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() || path == dirNode.Abs() || err != nil {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if !targetImpls(path) {
			return nil
		}
		logger.With("path", path).Info("Start handle impl")
		impl, err := NewBffImpl(types.NewBaseFile(path), sys)
		if err != nil {
			logger.Warnf("Fail to parse interface: %v", err.Error())
		} else if impl != nil {
			bff.Impls = append(bff.Impls, impl)
		}
		return nil
	})

	return bff, err
}

func (bffItem BffItem) ImplIter(iterFn func(impl *BffImpl) bool) {
	for _, impl := range bffItem.Impls {
		if !iterFn(impl) {
			return
		}
	}
}

func (bffItem BffItem) MatchImpl(fileNode types.FileNode) *BffImpl {
	if strings.HasPrefix(fileNode.Abs(), bffItem.DirNode.Abs()) {
		return nil
	}
	for _, bff := range bffItem.Impls {
		interfaceTree := bff.AstTree.Interface.(*BffInterfaceTree)
		if fileNode.Abs() == interfaceTree.FileNode.Abs() {
			return bff
		}
	}
	return nil
}

func (bffItem BffItem) MatchImplByName(implName string) *BffImpl {
	for _, bff := range bffItem.Impls {
		interfaceTree := bff.AstTree.Interface.(*BffInterfaceTree)
		if interfaceTree.FileNode.Name() == implName {
			return bff
		}
	}
	return nil
}

type BffImpl struct {
	Sys      *SystemDes `json:",omitempty" yaml:"-"`
	AstTree  *BffTree   `json:"-" yaml:"-"`
	FileNode types.FileNode
}

func NewBffImpl(fileNode types.FileNode, sysDes *SystemDes) (*BffImpl, error) {
	var bffImpl = new(BffImpl)
	var err error
	var bffTree *BffTree
	defer func() {
		if err1 := recover(); err1 != nil {
			err = errors.Format("%v", err1)
		}
	}()

	bffTree, err = ParseBffInterfaceRequestAstTree(fileNode)
	if err != nil {
		return nil, err
	}
	bffImpl.Sys = sysDes
	bffImpl.AstTree = bffTree
	bffImpl.FileNode = fileNode
	return bffImpl, nil
}

// Fuck，为聊兼容老项目
var skipImpls = func(_ string) bool {
	return false
}

func SkipImpls(fn func(_ string) bool) {
	skipImpls = fn
}

var targetImpls = func(_ string) bool {
	return true
}

func TargetImpls(fn func(_ string) bool) {
	targetImpls = fn
}

//
//type Services struct {
//	AbsDir       string        `json:",omitempty" yaml:"-"`
//	ServiceItems []ServiceItem `json:",omitempty" yaml:"names"`
//	Sys          *SystemDes    `json:",omitempty" yaml:"-"`
//}
//
//type ServiceItem struct {
//	Name   string        `json:",omitempty" yaml:"name,omitempty"`
//	AbsDir string        `json:",omitempty" yaml:"-"`
//	Impls  []ServiceImpl `json:",omitempty" yaml:"interfaces"`
//	Sys    *SystemDes    `json:",omitempty" yaml:"-"`
//}
//
//type ServiceImpl struct {
//	Name     string     `json:",omitempty" yaml:"name,omitempty"`
//	FileName string     `json:",omitempty" yaml:"filename,omitempty"`
//	AbsPath  string     `json:",omitempty" yaml:"-"`
//	Sys      *SystemDes `json:",omitempty" yaml:"-"`
//}
//
//type AssetsItem struct {
//	Name               string    `json:"name"`
//	CreateAt           time.Time `json:"created_at"`
//	BrowserDownloadUrl string    `json:"browser_download_url"`
//}
//
//type ApiRespItem struct {
//	Id       int          `json:"id"`
//	TagName  string       `json:"tag_name"`
//	CreateAt time.Time    `json:"created_at"`
//	Assets   []AssetsItem `json:"assets"`
//}
//
//func CurOsAsset(assets []AssetsItem) AssetsItem {
//	switch types.OsType() {
//	case types.MacOs:
//		return MacAsset(assets)
//	case types.LinuxOs:
//		return LinuxAsset(assets)
//	case types.WindowsOs:
//		return WindowsAsset(assets)
//	}
//
//	panic("unimplementable")
//}
//
//func MacAsset(assets []AssetsItem) AssetsItem {
//	for _, asset := range assets {
//		if strings.Contains(asset.Name, "mac") {
//			return asset
//		}
//	}
//	panic("unimplementable")
//}
//
//func LinuxAsset(assets []AssetsItem) AssetsItem {
//	for _, asset := range assets {
//		if strings.Contains(asset.Name, "linux") {
//			return asset
//		}
//	}
//	panic("unimplementable")
//}
//
//func WindowsAsset(assets []AssetsItem) AssetsItem {
//	for _, asset := range assets {
//		if strings.Contains(asset.Name, "windows") {
//			return asset
//		}
//	}
//	panic("unimplementable")
//}
