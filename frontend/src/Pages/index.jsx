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

    const handlerClick =  () => {
        console.log(value); // 
        if (value === ""){
            alert("No puedes enviar un comando vacío");
            return;
        } 
        Service.analisis(value)
        .then((res) => {
            setResponse(res.respuesta);
        })
        .catch((err) => {
            console.error(err);
        });
    }

    const handlerLimpiar = () => {
        if (value === ""){
            alert("No puedes limpiar un campo vacío");
            return;
        }
        changeText("");
        setResponse("");
    }

    const handleLoadClick = () => {
        const input = document.createElement('input');
        input.type = 'file';
        input.addEventListener('change', handleFileChange);
        input.click();
    }

    const handleFileChange = (e) => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.onload = (e) => {
            const text = e.target.result;
            changeText(text);
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

            <Consola text={"Consola de Salida"} handlerChange={setResponse} value={response} readOnly={true}/>
            </div>

        </>

    )
}

export default Index;
