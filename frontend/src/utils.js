function generateSuccessTextOutput(a,b){
    return `<span style='color: green'>${a}</span><span style='color: saddlebrown' >【`+ b + `】</span><br>`
}

function generateNormalTextOutput(a,b){
    return `<span style='color: ${b}'>${a}</span><br>`
}


export {
    generateSuccessTextOutput,
    generateNormalTextOutput
}