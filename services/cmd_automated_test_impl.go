package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/geekymedic/neon/logger"

	"github.com/geekymedic/neon-cli/templates"
	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/types/sysdes"

	"github.com/tidwall/pretty"
)

func (s *GenerateServer) GenerateAutomatedTest(ctx context.Context, arg *GenServerAutomatedTestArg) (_ *EmptyReply, err error) {
	// IntellijAutomated Test
	s.sys.Bffs.ImplIter(func(item *sysdes.BffItem, impl *sysdes.BffImpl) bool {
		var (
			bffInterfaceTree      = impl.AstTree.Interface.(*sysdes.BffInterfaceTree)
			bffReqTree            = impl.AstTree.Request.(*sysdes.BffRequestTree)
			intellijAutomatedTest templates.IntellijAutomatedProperty
			txt                   string
		)
		intellijAutomatedTest.Title = bffInterfaceTree.Annotation.Zh
		intellijAutomatedTest.URL = bffInterfaceTree.Annotation.URI
		if intellijAutomatedTest.URL == "" {
			intellijAutomatedTest.URL = fmt.Sprintf("/api/%s", PacketRouter(item.DirNode.Name(),
				bffInterfaceTree.FileNode.Name(), s.sys.Name))
		}

		var data bytes.Buffer
		if err = json.NewEncoder(&data).Encode(InjectAstTree(bffReqTree.TopNode)); err != nil {
			return false
		}
		intellijAutomatedTest.Data = fmt.Sprintf("%s", pretty.Pretty(data.Bytes()))
		txt, err = templates.ParseTemplate(templates.IntellijAutomatedTemplate, intellijAutomatedTest)
		types.AssertNil(err)

		filename := arg.Out.Append(fmt.Sprintf("%s_%s.http", item.DirNode.Name(), impl.FileNode.Name())).(types.DirNode).Abs()
		mkFp := types.NewBaseFile(filename)
		err = mkFp.Create(os.O_CREATE|os.O_WRONLY|os.O_EXCL, types.DefPerm)
		if err != nil {
			logger.With("path", mkFp.Abs()).Info("Skipped")
			return true
		}
		_, err = mkFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(mkFp.Sync())
		types.AssertNil(mkFp.Close())
		return true
	})

	// 生成http-client-private.env.json

	return
}
