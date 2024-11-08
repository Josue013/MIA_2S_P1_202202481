import React from "react";
import Editor from "@monaco-editor/react";
import "../Styles/Editor.css";

function Consola(props){
  const handlerChangeEditor = (newValue, event) => {
    props.handlerChange(newValue);
  }

  return(
    <>
    <div className="mb-3" editor-react>
            <label
            htmlFor="exampleFormControlTextarea1"
            style={{ fontSize: "30px" }}
            >
                {props.text}
            </label>
            <Editor
            height="50vh"
            width="100vh"
            theme="vs-dark"
            defaultLanguage="bash"
            value = {props.value}
            onChange={handlerChangeEditor}
            options={{
                selectOnLineNumbers: true,
                automaticLayout: true,
                lineNumbers: "on",
                readOnly: props.readOnly,
            }}
            />
        </div>
    </>
  );
}

export default Consola;