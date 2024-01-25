const formbtn = document.querySelector('.btn');
// const inputs = document.querySelectorAll('.form-inputs');

const formValues = document.querySelector('form');

formbtn.addEventListener('click', async (e) => {
    e.preventDefault();
    const fv = new FormData(formValues);
    const obj = Object.fromEntries(fv);
    const jsonData = JSON.stringify(obj);
    const datafetch = await fetch('/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
          },
        body: jsonData
     });

    let result = await datafetch.json();

    // Return -> main page
        // window.location.assign('http://127.0.0.1:5000/');
    

    console.log(result)
    console.log(jsonData)
})