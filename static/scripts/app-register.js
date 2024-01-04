const formbtn = document.querySelector('.btn');
// const inputs = document.querySelectorAll('.form-inputs');

const formValues = document.querySelector('form');

formbtn.addEventListener('click', async (e) => {
    e.preventDefault();
    const fv = new FormData(formValues);
    const obj = Object.fromEntries(fv);
    const jsonData = JSON.stringify(obj);
    console.log("sending JSON:", jsonData)
    const datafetch = await fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
          },
        body: jsonData
     });

    let result = await datafetch.json();

    const h1Element = document.createElement('h1');
    if (result.status_code === 200) {
        console.log('Registration ok')
        h1Element.textContent = 'Registration ok';
    } else {
        console.log(`Registration error with. Details: ${result.message}!`);
        h1Element.textContent = `Registration error with. Details: ${result.message}!`;
    }
    document.body.appendChild(h1Element);

    // Return -> main page
        // window.location.assign('http://127.0.0.1:5000/');


    // console.log(result)
    // console.log(jsonData)
})