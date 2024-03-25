"use strict"
const formbtn = document.querySelector('.btn');
const formInputs = document.querySelectorAll('.form-inputs');
const formValues = document.querySelector('form');
const inputsContainer = document.querySelector('.inputs-container');

function keyDownEvents(){
    formbtn.disabled = true;
    for(let elem of formInputs){
        elem.addEventListener('input', () => {
            if(elem.value.length < 1){
                formbtn.disabled = true
                return
            } 
                formbtn.disabled = false
        })
        
    }
 
}

keyDownEvents()

formbtn.addEventListener('click', async (e) => {
    e.preventDefault();
    loader.style.display = 'block'
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
            let user_succ = document.createElement('small');
                    user_succ.classList.add('big-user-succ');
                    user_succ.innerText = `Success!`;
                    // remove err if already exist
                    if(inputsContainer.firstChild.classList?.contains('big-user-succ')
                     || inputsContainer.firstChild.classList?.contains('big-user-err')){
                        inputsContainer.firstChild.remove();
                    }
                    inputsContainer.prepend(user_succ);
                    setTimeout(() => {
                        user_succ.remove();
                    }, 5000);
            console.log('successful login');
            location.reload();
            return;
        case 302:
        case 303:
            console.log('repetitive login attempt was redirected');
            location.reload();
            return 
        case 400:
            console.log(datafetch.statusText, 'Empty values.')
            loader.style.display = 'none'
            return 
        case 401:
            
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
            console.log(datafetch.statusText, 'Invalid username or password.')
            loader.style.display = 'none'
            return 
        case 500:
            console.log(datafetch.statusText, '')
            loader.style.display = 'none'
            return 
        default:
            console.log(`unexpected status received: ${datafetch.status}`)
            loader.style.display = 'none'
    }
    let result = await datafetch.json();
    console.log(result)
    console.log(jsonData)

    
})


