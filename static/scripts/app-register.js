"use strict"
const formbtn = document.querySelector('.btn');
const formValues = document.querySelector('form');
const formInputs = document.querySelectorAll('.form-inputs');
const inputsContainer = document.querySelector('.inputs-container');

function keyDownEvents(){
    if(formbtn.getAttribute('aria-disabled') === 'true'){
        return
    }

    const REGEXP = regExpDelivery();

    for(let elem of formInputs){
        elem.addEventListener('input', () => {
            
            let user_err = document.createElement('small');
            if(elem.nextElementSibling?.classList.contains('user-err')){
                elem.nextElementSibling.remove();
            }

            if(!REGEXP.get(elem.name).test(elem.value)){
                user_err.classList.add('user-err');
                user_err.innerText = keydownMessages(elem).get(elem.name);
                elem.style.border = '2px solid #b90909'
                elem.after(user_err);
                formbtn.setAttribute('aria-disabled', 'true')
            } else if(REGEXP.get(elem.name).test(elem.value)){
                user_err.remove();
                elem.style.border = '2px solid #dddddd';
                elem.addEventListener('focus', () => elem.style.border = '2px solid #0941b9');
                elem.addEventListener('blur', () => elem.style.border = '2px solid #dddddd')
                formbtn.setAttribute('aria-disabled', 'false')
            }
        })
    }
 
}

keyDownEvents()

function keydownMessages(inputName){
    let errorMessages = new Map();
    errorMessages.set('username', `${inputName.name.slice(0, 1).toUpperCase() + inputName.name.slice(1)} must be in range from 5 to 20 characters`)
                .set('email', `${inputName.name.slice(0, 1).toUpperCase() + inputName.name.slice(1)} should look like this : example@gmail.com`)
                .set('password', `Make sure ${inputName.name.slice(0, 1).toUpperCase() + inputName.name.slice(1)} it\'s at least 8 characters \n including: a lowercase letter, a uppercase letter, \n a number and special character`)

    return errorMessages
}




formbtn.addEventListener('click', async (e) => {
    e.preventDefault();

    if(formbtn.getAttribute('aria-disabled') === 'true'){
        return
    }

    const REGEXP = regExpDelivery();
    
    const formInputs = document.querySelectorAll('.form-inputs');
    const userNameInput = document.querySelector('#username');

    let validation = preValidation(formInputs, REGEXP);
    if(!validation.every(el => el === false)) {
        console.log('no validation happening');
        return;
    }

    // Validation & comms to user
    formbtn.setAttribute('aria-disabled', 'true')
    loader.style.visibility = 'visible'
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

        const result = await datafetch.json();

        if(datafetch.ok){
            console.log('successful registration');
            let user_succ = document.createElement('small');
            user_succ.classList.add('big-user-succ');
            user_succ.innerText = `You Successfully Register!`;
            // remove err if already exist
            if(inputsContainer.firstChild.classList?.contains('big-user-succ')
             || inputsContainer.firstChild.classList?.contains('big-user-err')){
                inputsContainer.firstChild.remove();
            }
            inputsContainer.prepend(user_succ);
            setTimeout(() => {
                user_succ.remove();
            }, 5000);
            window.location.href = '/login';
            return;
        }
            console.log(result.message)
            loader.style.visibility = 'hidden'
            let user_err = document.createElement('small');
            user_err.classList.add('user-err');
            user_err.innerText = `${result.message}`;
            userNameInput.style.border = '2px solid #b90909';
            // remove err if already exist
            if(userNameInput.nextElementSibling.classList.contains('user-err')){
                userNameInput.nextElementSibling.remove();
            }
            userNameInput.after(user_err);
            setTimeout(() => {
                user_err.remove();
                userNameInput.style.border = '2px solid #dddddd';
                userNameInput.addEventListener('focus', () => userNameInput.style.border = '2px solid #0941b9');
                userNameInput.addEventListener('blur', () => userNameInput.style.border = '2px solid #dddddd');
                formbtn.setAttribute('aria-disabled', 'false')
            }, 5000);
            return
})


// Regexp Logic
function regExpDelivery(){
    const REGEXP = new Map();
    REGEXP.set('username', /^(?=[\w.-]{5,20}$)(?:[\d_.-]*[a-zA-Z]){3}[\w.-]*$/)
        // ^ Start of string
        // (?=[\w.-]{5,20}$) Assert 5-20 of the allowed characters
        // (?:[\d_.-]*[a-zA-Z]){3} Match 3 times a-zA-Z followed by optional digits _ . -
        // [\w.-]* Match optional word chars or . or -
        // $ End of string
    .set('email', /^[a-z0-9]+@[a-z]+\.[a-z]{2,3}$/)
        //  ^ - It is the start of the string.
        // [a-z0-9]+ - Any character between a to z and 0 to 9 at the start of the string.
        // @ - The string should contains ‘@’ character after some alphanumeric characters.
        // [a-z]+ - At least one character between a to z after the ‘@’ character in the string.
        // \. – Email should contain the dot followed by some characters followed by the ‘@’ character
        // [a-z]{2,3}$ - It should contain two or three alphabetical characters at the end of the string. The ‘$’ represents the end of the string.
    .set('password', /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#$%^&*]).{8,72}$/)
        // old one : /^(?=.*[A-Z])(?=.*[0-9])(?=.*[a-z])[a-zA-Z0-9]{6,}$/
        // At least 8 characters long
        // contains a lowercase letter
        // contains an uppercase letter
        // contains a digit
        // contains a special character
    return REGEXP;
}

// prevalidation logic
function preValidation(inputsData, REGEXP){
    let arr = [];

    for(let elem of inputsData){
        
        arr.push(!REGEXP.get(elem.name).test(elem.value))
        if(!REGEXP.get(elem.name).test(elem.value)){
            let user_err = document.createElement('small');
            user_err.classList.add('user-err');
            user_err.innerText = `Invalid ${elem.name}`;
            elem.style.border = '2px solid #b90909'
            elem.after(user_err);
            formbtn.setAttribute('aria-disabled', 'true')


                setTimeout(() => {
                    user_err.remove();
                    elem.style.border = '2px solid #dddddd';
                    elem.addEventListener('focus', () => elem.style.border = '2px solid #0941b9');
                    elem.addEventListener('blur', () => elem.style.border = '2px solid #dddddd')
                    formbtn.setAttribute('aria-disabled', 'false')
                }, 5000);
        }   
    }
    return arr
}