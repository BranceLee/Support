import React, { useEffect, useState } from 'react';
import Quill from 'quill';
import axios from 'axios';
import qs from 'qs';

const baseURL = 'http://localhost:8000';

const saveContent = payload => {
  const option = {
    method: 'post',
    baseURL: baseURL,
    url: '/api/blog',
    headers: { 'content-type': 'application/x-www-form-urlencoded' },
    data: qs.stringify(payload),
  };
  return axios(option)
    .then(res => res.data)
    .catch(err => console.log(err));
};

const toolbarOptions = [
  ['bold', 'italic', 'underline', 'strike'],
  ['blockquote', 'code-block'],
  [{ header: [1, 2, 3, 4, 5, 6, false] }],
  [{ align: [] }],
];

function Index() {
  const [editorContent, setEditorContent] = useState('');

  let editor;
  useEffect(() => {
    editor = new Quill('#editor', {
      modules: { toolbar: toolbarOptions },
      theme: 'snow',
    });
  });

  const changeContent = () => {
    const contentWrap = document.getElementById('content');
    contentWrap.innerHTML = editor.root.innerHTML;
    // setEditorContent(editor.root.innerHTML);
   };

  const sendContentToServer = async () => {
    const payload = {
      title: 'hello',
      content: editor.root.innerHTML,
    };
    const data = await saveContent(payload);
  };

  return (
    <div>
      <div id="toolbar" />
      <div id="editor">
        <p>Hello World!</p>
      </div>
      <button onClick={() => changeContent()}>Preview</button>
      <button onClick={() => sendContentToServer()}>Save</button>
      <div className="ql-snow">
        <div className="ql-editor">
          <div id="content" />
        </div>
      </div>
    </div>
  );
}

export default Index;
