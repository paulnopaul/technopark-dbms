// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package domain

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain(in *jlexer.Lexer, out *Vote) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "nickname":
			out.Nickname = string(in.String())
		case "voice":
			out.Voice = int32(in.Int32())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain(out *jwriter.Writer, in Vote) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"nickname\":"
		out.RawString(prefix[1:])
		out.String(string(in.Nickname))
	}
	{
		const prefix string = ",\"voice\":"
		out.RawString(prefix)
		out.Int32(int32(in.Voice))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Vote) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Vote) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Vote) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Vote) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain(l, v)
}
func easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain1(in *jlexer.Lexer, out *UserArray) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(UserArray, 0, 1)
			} else {
				*out = UserArray{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 User
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain1(out *jwriter.Writer, in UserArray) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v UserArray) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserArray) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserArray) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserArray) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain1(l, v)
}
func easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain2(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "nickname":
			out.Nickname = string(in.String())
		case "fullname":
			out.Fullname = string(in.String())
		case "about":
			out.About = string(in.String())
		case "email":
			out.Email = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain2(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Nickname != "" {
		const prefix string = ",\"nickname\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Nickname))
	}
	if in.Fullname != "" {
		const prefix string = ",\"fullname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Fullname))
	}
	if in.About != "" {
		const prefix string = ",\"about\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.About))
	}
	if in.Email != "" {
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain2(l, v)
}
func easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain3(in *jlexer.Lexer, out *ThreadArray) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(ThreadArray, 0, 0)
			} else {
				*out = ThreadArray{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v4 Thread
			(v4).UnmarshalEasyJSON(in)
			*out = append(*out, v4)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain3(out *jwriter.Writer, in ThreadArray) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v5, v6 := range in {
			if v5 > 0 {
				out.RawByte(',')
			}
			(v6).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v ThreadArray) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ThreadArray) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ThreadArray) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ThreadArray) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain3(l, v)
}
func easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain4(in *jlexer.Lexer, out *Thread) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int32(in.Int32())
		case "title":
			out.Title = string(in.String())
		case "author":
			out.Author = string(in.String())
		case "forum":
			out.Forum = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "votes":
			out.Votes = int32(in.Int32())
		case "slug":
			out.Slug = string(in.String())
		case "created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Created).UnmarshalJSON(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain4(out *jwriter.Writer, in Thread) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int32(int32(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix)
		out.String(string(in.Author))
	}
	{
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.String(string(in.Forum))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	if in.Votes != 0 {
		const prefix string = ",\"votes\":"
		out.RawString(prefix)
		out.Int32(int32(in.Votes))
	}
	if in.Slug != "" {
		const prefix string = ",\"slug\":"
		out.RawString(prefix)
		out.String(string(in.Slug))
	}
	if true {
		const prefix string = ",\"created\":"
		out.RawString(prefix)
		out.Raw((in.Created).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Thread) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Thread) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Thread) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Thread) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain4(l, v)
}
func easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain5(in *jlexer.Lexer, out *Service) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "user":
			out.User = int32(in.Int32())
		case "forum":
			out.Forum = int32(in.Int32())
		case "thread":
			out.Thread = int32(in.Int32())
		case "post":
			out.Post = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain5(out *jwriter.Writer, in Service) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix[1:])
		out.Int32(int32(in.User))
	}
	{
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.Int32(int32(in.Forum))
	}
	{
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		out.Int32(int32(in.Thread))
	}
	{
		const prefix string = ",\"post\":"
		out.RawString(prefix)
		out.Int64(int64(in.Post))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Service) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Service) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Service) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Service) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain5(l, v)
}
func easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain6(in *jlexer.Lexer, out *PostFull) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "post":
			if in.IsNull() {
				in.Skip()
				out.Post = nil
			} else {
				if out.Post == nil {
					out.Post = new(Post)
				}
				(*out.Post).UnmarshalEasyJSON(in)
			}
		case "forum":
			if in.IsNull() {
				in.Skip()
				out.Forum = nil
			} else {
				if out.Forum == nil {
					out.Forum = new(Forum)
				}
				(*out.Forum).UnmarshalEasyJSON(in)
			}
		case "thread":
			if in.IsNull() {
				in.Skip()
				out.Thread = nil
			} else {
				if out.Thread == nil {
					out.Thread = new(Thread)
				}
				(*out.Thread).UnmarshalEasyJSON(in)
			}
		case "author":
			if in.IsNull() {
				in.Skip()
				out.User = nil
			} else {
				if out.User == nil {
					out.User = new(User)
				}
				(*out.User).UnmarshalEasyJSON(in)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain6(out *jwriter.Writer, in PostFull) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"post\":"
		out.RawString(prefix[1:])
		if in.Post == nil {
			out.RawString("null")
		} else {
			(*in.Post).MarshalEasyJSON(out)
		}
	}
	if in.Forum != nil {
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		(*in.Forum).MarshalEasyJSON(out)
	}
	if in.Thread != nil {
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		(*in.Thread).MarshalEasyJSON(out)
	}
	if in.User != nil {
		const prefix string = ",\"author\":"
		out.RawString(prefix)
		(*in.User).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostFull) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostFull) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostFull) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostFull) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain6(l, v)
}
func easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain7(in *jlexer.Lexer, out *PostArray) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(PostArray, 0, 0)
			} else {
				*out = PostArray{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v7 Post
			(v7).UnmarshalEasyJSON(in)
			*out = append(*out, v7)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain7(out *jwriter.Writer, in PostArray) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v8, v9 := range in {
			if v8 > 0 {
				out.RawByte(',')
			}
			(v9).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v PostArray) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostArray) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostArray) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostArray) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain7(l, v)
}
func easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain8(in *jlexer.Lexer, out *Post) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "parent":
			out.Parent = int64(in.Int64())
		case "author":
			out.Author = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "isEdited":
			out.IsEdited = bool(in.Bool())
		case "forum":
			out.Forum = string(in.String())
		case "thread":
			out.Thread = int32(in.Int32())
		case "created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Created).UnmarshalJSON(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain8(out *jwriter.Writer, in Post) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	if in.Parent != 0 {
		const prefix string = ",\"parent\":"
		out.RawString(prefix)
		out.Int64(int64(in.Parent))
	}
	if in.Author != "" {
		const prefix string = ",\"author\":"
		out.RawString(prefix)
		out.String(string(in.Author))
	}
	if in.Message != "" {
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	if in.IsEdited {
		const prefix string = ",\"isEdited\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsEdited))
	}
	if in.Forum != "" {
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.String(string(in.Forum))
	}
	if in.Thread != 0 {
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		out.Int32(int32(in.Thread))
	}
	if true {
		const prefix string = ",\"created\":"
		out.RawString(prefix)
		out.Raw((in.Created).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Post) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Post) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Post) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Post) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain8(l, v)
}
func easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain9(in *jlexer.Lexer, out *JSONMessageType) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "message":
			out.Message = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain9(out *jwriter.Writer, in JSONMessageType) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix[1:])
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v JSONMessageType) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v JSONMessageType) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *JSONMessageType) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *JSONMessageType) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain9(l, v)
}
func easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain10(in *jlexer.Lexer, out *Forum) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "title":
			out.Title = string(in.String())
		case "user":
			out.User = string(in.String())
		case "slug":
			out.Slug = string(in.String())
		case "posts":
			out.Posts = int64(in.Int64())
		case "threads":
			out.Threads = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain10(out *jwriter.Writer, in Forum) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix[1:])
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix)
		out.String(string(in.User))
	}
	{
		const prefix string = ",\"slug\":"
		out.RawString(prefix)
		out.String(string(in.Slug))
	}
	if in.Posts != 0 {
		const prefix string = ",\"posts\":"
		out.RawString(prefix)
		out.Int64(int64(in.Posts))
	}
	if in.Threads != 0 {
		const prefix string = ",\"threads\":"
		out.RawString(prefix)
		out.Int64(int64(in.Threads))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Forum) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Forum) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeTechnoparkDbmsInternalPkgDomain10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Forum) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Forum) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeTechnoparkDbmsInternalPkgDomain10(l, v)
}
