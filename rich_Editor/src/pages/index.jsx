import React, {useEffect, useState} from 'react'
import Quill from 'quill'
import axios from 'axios'
import qs from 'qs'

const baseURL = 'http://localhost:8000'

const saveContent=(payload)=>{
    const option = {
        method:'post',
        baseURL:baseURL,
        url:"/api/blog",
        headers:{'content-type':'application/x-www-form-urlencoded'},
        data:qs.stringify(payload)
    }
    return axios(option).then(res=>res.data).catch(err=>console.log(err))
}
function Index(){
    const [editorContent, setEditorContent]= useState('')

    let editor
    useEffect(()=>{
        editor = new Quill('#editor',{
            modules:{toolbar:'#toolbar'},
            theme:'snow',
        });
        
    })
    const changeContent=()=>{
        const contentWrap=document.getElementById('content');
        contentWrap.innerHTML = editor.root.innerHTML;
    }

    const sendContentToServer=async()=>{
        const payload = {
            title:      'hello',
            content:    editorContent
        }
        const data=await saveContent(payload)
    }

    return (
        <div>
            <div id='toolbar'>
                <button class="ql-bold">Bold</button>
                <button class="ql-italic">Italic</button>
            </div>
            <div id='editor'>
                 <p>Hello World!</p>
            </div>
            <button onClick={()=>changeContent()}>Active</button>
            <div id='content'/>
            <button onClick={()=>sendContentToServer()}>Save</button>
        </div>
    )
}

export default Index