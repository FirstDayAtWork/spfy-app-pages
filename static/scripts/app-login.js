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
    switch (datafetch.status) {
        case 200:
            console.log('successful login');
            location.reload();
            return;
        case 302:
        case 303:
            console.log('repetitive login attempt was redirected');
            location.reload();
            return;
        default:
            console.log(`unexpected status received: ${datafetch.status}`)
    }
    let result = await datafetch.json();
    console.log(result)
    console.log(jsonData)

    
})