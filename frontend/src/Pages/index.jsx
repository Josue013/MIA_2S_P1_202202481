import React from "react";
import NavBar from "../components/NavBar";
import Consola from "../components/Editor";
import Service from "../Services/Service";
import '../Styles/NavBar.css'

function Index() {
    const [value, setValue] = React.useState("");
    const [response, setResponse] = React.useState("");

    const changeText = (text) => {
        setValue(text);
    }

    const handlerClick = () => {
        if (value === "") {
            alert("No hay nada que analizar");
            return;
        }
        Service.analisis(value)
            .then((response) => {
                setResponse(response);
            })
    }

    const handlerLimpiar = () => {
        if (value === "") {
            alert("No hay nada que limpiar");
            return;
        }
        changeText("");
        setResponse("");
    }

    const handleLoadClick = () => {
        const input = document.createElement("input");
        input.type = "file";
        input.addEventListener("change", handleFileChange);
        input.click();
    }

    const handleFileChange = (e) => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.onload = (e) => {
            const text = e.target.result;
            setValue(text);
        }
        reader.readAsText(file);
    }

    return (
        <>
            <div className="item-0">
                <div className="container-NavB">
                    <div className="item-1-NavB">
                        <button className="button-74" role="button" onClick={handleLoadClick}>Abrir Archivo</button>
                    </div>
                    <div className="item-2-NavB">
                        <button className="button-74" role="button" onClick={handlerClick}>Ejecutar</button>
                    </div>
                    <div className="item-3-NavB">
                        <button className="button-74" role="button" onClick={handleLoadClick}>Reportes</button>
                    </div>
                </div>
            </div>

            <h1>Proyecto 1 - MIA - 202202481</h1>
            
            <div class="container">
                <button type="button" class="btn btn-primary" onClick={handlerLimpiar}>Limpiar</button>
            </div>

            <div class="Consola">
            <Consola text={"Consola de Entrada"} handlerChange={changeText} value={value} />

            <Consola text={"Consola de Salida"} handlerChange={changeText} value={response} readOnly />
            </div>

        </>

    )
}

export default Index;
