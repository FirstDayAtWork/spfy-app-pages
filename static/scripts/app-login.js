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

    const result = await datafetch.json();
    
    if(datafetch.ok){
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
    }
        let user_err = document.createElement('small');
        user_err.classList.add('big-user-err');
        user_err.innerText = `${result.message}`;
        // remove err if already exist
        if(inputsContainer.firstChild.classList?.contains('big-user-err')){
            inputsContainer.firstChild.remove();
        }
        inputsContainer.prepend(user_err);
        setTimeout(() => {
            user_err.remove();
        }, 5000);
        console.log(result.message)
        loader.style.display = 'none'
        return 
})


