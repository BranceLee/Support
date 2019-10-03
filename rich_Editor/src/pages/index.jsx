import React, {useEffect, useState} from 'react'
import Quill from 'quill'

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
        </div>
    )
}

export default Index