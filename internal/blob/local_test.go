package blob

// func TestLocalFS(t *testing.T) {
// 	vol := "test/data/blob"
// 	defer os.Remove("test")
// 	defer os.Remove("test/data")
// 	defer os.RemoveAll(vol)

// 	filename := "myfile.txt"
// 	data1 := bytes.Repeat([]byte{'X'}, 4*humanize.KiByte)
// 	data2 := append(data1, bytes.Repeat([]byte{'Y'}, 4*humanize.KiByte)...)
// 	data3 := append(data2, bytes.Repeat([]byte{'Z'}, 4*humanize.KiByte)...)
// 	data4 := append(data3, bytes.Repeat([]byte{'A'}, 4*humanize.KiByte)...)
// 	data5 := append(data4, bytes.Repeat([]byte{'B'}, 4*humanize.KiByte)...)
// 	src := bytes.NewReader(data5)
// 	if src.Len() == 0 {
// 		t.Fatal("src size cannot zero : ", src.Len())
// 	}
// 	lfs, err := newLocalFS(vol)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	ref := &Ref{
// 		ID:          uuid.Nil,
// 		UserID:      uuid.Nil,
// 		ObjectName:  filename,
// 		ContentType: "",
// 		Size:        0,
// 		Sum:         "",
// 		CreatedAt:   time.Time{},
// 	}

// 	err = lfs.Put(context.TODO(), ref, src)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if ref.Size != uint64(20*humanize.KiByte) {
// 		t.Fatal("got: ", ref.Size)
// 	}

// 	rd, _, err := lfs.Get(context.Background(), ref)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	dst := bytes.Buffer{}
// 	n, err := io.Copy(&dst, rd)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if n != int64(ref.Size) {
// 		t.Fatal("got ", ref.Size, "expect", n)
// 	}

// 	if !bytes.Equal(data5, dst.Bytes()) {
// 		t.Fatal("not equal")
// 	}

// 	if err := rd.Close(); err != nil {
// 		t.Fatal(err)
// 	}

// }

// func Benchmark_LocalFS_Put(b *testing.B) {
// 	vol := "test/data/blob"
// 	filename := "myfile.txt"
// 	data1 := bytes.Repeat([]byte{'X'}, 4*humanize.KiByte)
// 	data2 := append(data1, bytes.Repeat([]byte{'Y'}, 4*humanize.KiByte)...)
// 	data3 := append(data2, bytes.Repeat([]byte{'Z'}, 4*humanize.KiByte)...)
// 	data4 := append(data3, bytes.Repeat([]byte{'A'}, 4*humanize.KiByte)...)
// 	data5 := append(data4, bytes.Repeat([]byte{'B'}, 4*humanize.KiByte)...)
// 	src := bytes.Reader{}

// 	lfs, _ := newLocalFS(vol)
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		b.StopTimer()
// 		src.Reset(data5)
// 		ref := Ref{
// 			ID:          uuid.Nil,
// 			UserID:      uuid.Nil,
// 			ObjectName:  fmt.Sprintf("%v%d", filename, b.N),
// 			ContentType: "",
// 			Size:        0,
// 			Sum:         "",
// 			CreatedAt:   time.Time{},
// 		}
// 		b.StartTimer()

// 		//start measuring
// 		if err := lfs.Put(context.Background(), &ref, &src); err != nil {
// 			b.Fatal(err)
// 		}
// 		//end measuring

// 	}

// 	b.StopTimer()
// 	os.RemoveAll(vol)
// 	b.StartTimer()

// }

// func Benchmark_LocalFS_Get(b *testing.B) {
// 	vol := "test/data/blob"
// 	filename := "myfile.txt"
// 	data1 := bytes.Repeat([]byte{'X'}, 4*humanize.KiByte)
// 	data2 := append(data1, bytes.Repeat([]byte{'Y'}, 4*humanize.KiByte)...)
// 	data3 := append(data2, bytes.Repeat([]byte{'Z'}, 4*humanize.KiByte)...)
// 	data4 := append(data3, bytes.Repeat([]byte{'A'}, 4*humanize.KiByte)...)
// 	data5 := append(data4, bytes.Repeat([]byte{'B'}, 4*humanize.KiByte)...)
// 	src := bytes.Reader{}
// 	src.Reset(data5)

// 	lfs, _ := newLocalFS(vol)
// 	ref := Ref{
// 		ID:          uuid.Nil,
// 		UserID:      uuid.Nil,
// 		ObjectName:  filename,
// 		ContentType: "",
// 		Size:        0,
// 		Sum:         "",
// 		CreatedAt:   time.Time{},
// 	}
// 	if err := lfs.Put(context.Background(), &ref, &src); err != nil {
// 		b.Fatal(err)
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		r, _, err := lfs.Get(context.Background(), &ref)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 		if _, err := io.Copy(io.Discard, r); err != nil {
// 			b.Fatal(err)
// 		}
// 		r.Close()
// 	}

// 	b.StopTimer()
// 	os.RemoveAll(vol)
// 	b.StartTimer()
// }
