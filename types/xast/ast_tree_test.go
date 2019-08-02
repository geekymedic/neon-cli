package xast

//var astFileBinary = `
//	package foo
//
//	type Person struct {
//		Name        string
//		Address     Location
//		Degree      map[string]int
//		LinkAddress map[string]Point
//		HairColor   []string
//		LegHeight   []Height
//	}
//
//	type Location struct {
//		Stress  string
//		Pointer Point
//	}
//
//	type Point struct {
//		Latitude  int
//		Longitude int
//	}
//
//	type Height struct {
//		Value int
//	}
//`

//func buildTree() TopNode {
//	root := TopNode{
//		TypeName: "Person",
//		LeavesNodes: map[string]*LeafNode{
//			"Name": {
//				TypeName: "string",
//			},
//		},
//		ExtraNodes: map[string]*ExtraNode{
//			"Address": {
//				TypeName: "Location",
//				LeavesNodes: map[string]*LeafNode{
//					"Stress": {
//						TypeName: "string",
//					},
//				},
//				ExtraNodes: map[string]*ExtraNode{
//					"Pointer": {
//						TypeName: "Point",
//						LeavesNodes: map[string]*LeafNode{
//							"Latitude": {TypeName: "int"}, "Longitude": {TypeName: "int"},
//						},
//					},
//				},
//			},
//
//			"Degree": {
//				TypeName: "map",
//				LeavesNodes: map[string]*LeafNode{VirMap + "-int" + "-string": {
//					TypeName: "int",
//				}},
//			},
//
//			// map[string]Point
//			"LinkAddress": {
//				TypeName: "map",
//				ExtraNodes: map[string]*ExtraNode{
//					VirMap + "-" + "Point": {
//						TypeName: "Point",
//						LeavesNodes: map[string]*LeafNode{
//							"Latitude": {
//								TypeName: "int",
//							},
//							"Longitude": {
//								TypeName: "int",
//							},
//						},
//					},
//				},
//			},
//
//			"HairColor": {
//				TypeName: "array",
//				LeavesNodes: map[string]*LeafNode{
//					"Height": {
//						TypeName: VirArray + "-" + "Height",
//					},
//				},
//			},
//
//			"LegHeight": {
//				TypeName: "array",
//				ExtraNodes: map[string]*ExtraNode{
//					VirArray + "-" + "Height": {
//						TypeName: "Height",
//						ExtraNodes: map[string]*ExtraNode{
//							"Value": {
//								TypeName: "string",
//							},
//						},
//					},
//				},
//			},
//		},
//	}
//
//	return root
//}

//func TestTopNode_FindNode(t *testing.T) {
//	t.Run("Should be found node", func(t *testing.T) {
//		var args = []struct {
//			Param  string
//			Expect bool
//		}{
//			{"Person", true},
//			{"Person.Name", true},
//		}
//		for _, arg := range args {
//			root := buildTree()
//			_, actual := root.FindNode(arg.Param)
//			if actual != arg.Expect {
//				t.Fatalf("expect %v, got %v", arg.Expect, actual)
//			}
//		}
//	})
//
//	t.Run("Should be insert node successfully", func(t *testing.T) {
//		var args = []struct {
//			ParentNode   string
//			NewExtraNode ExtraNode
//			NewVar       string
//			Expect       error
//		}{
//			{"Person.Address", ExtraNode{
//				"NewTest",
//				nil,
//				nil,
//				nil,
//				"Person.Address.NewVarName"}, "NewVarName", nil},
//		}
//
//		for _, arg := range args {
//			root := buildTree1()
//			_, actual := root.FindNode(arg.ParentNode)
//			if !actual {
//				t.Fatalf("expect %v, actual %v", arg.Expect, actual)
//			}
//			if actual := root.AfterInsertExtraNode(arg.ParentNode, arg.NewVar, arg.NewExtraNode); arg.Expect != actual {
//				t.Fatalf("expect %v, actual %v", arg.Expect, actual)
//			}
//			root.ReBuildWalkPath()
//			json.NewEncoder(os.Stdout).Encode(root)
//		}
//	})
//
//	t.Run("Should be not found node", func(t *testing.T) {
//		root := buildTree1()
//		root.FindNode("ddddd")
//		json.NewEncoder(os.Stdout).Encode(root)
//	})
//}
