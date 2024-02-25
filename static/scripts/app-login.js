const formbtn = document.querySelector('.btn');
const inputs = document.querySelectorAll('.form-inputs');

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
        case 401:
            const inputsContainer = document.querySelector('.inputs-container');
            let user_err = document.createElement('small');
                    user_err.classList.add('big-user-err');
                    user_err.innerText = `Invalid username or password.`;
                    // remove err if already exist
                    if(inputsContainer.firstChild.classList?.contains('big-user-err')){
                        inputsContainer.firstChild.remove();
                    }
                    inputsContainer.prepend(user_err);

                    setTimeout(() => {
                        user_err.remove();
                    }, 5000);
            console.log(datafetch.statusText, 'Invalid USERNAME or PASSWORD.')
            return
        default:
            console.log(`unexpected status received: ${datafetch.status}`)
    }
    let result = await datafetch.json();
    console.log(result)
    console.log(jsonData)

    
})