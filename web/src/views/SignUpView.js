import AbstractView from "./AbstractView.js";
import fetcher from "../services/Fetcher.js";
import router from "../index.js"

const signUp = async (body) => {
    const path = `/api/users/sign-up`
    const data = await fetcher.post(path, body)
    if (data && data.error) {
        drawError(data.error)
        return
    }
    router.navigateTo("/sign-in")
}

const drawError = (err) => {
    const errorMessage = document.getElementById("error-message")
    errorMessage.innerText = err
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Sign up");
    }

    async getHtml() {
        return `
            <form id="sign-up-form" onsubmit="return false;">
                <label for="username">Pseudo</label>
                <input type="text" id="username" placeholder="Pseudo" required minlength="2" maxlength="64" pattern="^(?![_.])(?!.*[_.-]{2})[a-zA-Z0-9._-]+(?<![_.-])$" title="Username should only contain alphanumerical and '.', '_', '-' symbols, no symbol at the beginnig and at the end, no alternation of special characters"> </br> </br>

                <label for="first-name">Prénom</label>
                <input type="text" id="first-name" placeholder="Prénom" required minlength="2" maxlength="64" pattern="[a-zA-Z]+$" title="First name should only contain latin letters">  <br> <br>

                <label for="last-name">Nom</label>
                <input type="text" id="last-name" placeholder="Nom" required minlength="2" maxlength="64"  pattern="[a-zA-Z]+$" title="Last name should only contain latin letters"> <br> <br>

                <label for="age">Age</label>
                <input type="number" id="age" placeholder="Age" required min="12" max="110"> <br> <br>

                <p>Genre</p>
                <input type="radio" name="gender" id="gender-male" value="1" required>
                <label for="gender-male">Homme</label>
                <input type="radio" name="gender" id="gender-female" value="2">
                <label for="gender-female">Femme</label><br> <br>

                <label for="email">E-mail</label>
                <input type="email" id="email" placeholder="E-mail" required maxlength="64">  <br><br>
                
                <label for="password">Mot de passe</label>
                <input type="password" id="password" placeholder="Mot de passe" minlength=7 maxlength="64" required pattern="(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#~$%^&*()+|_]).{7,}" title="Password must contain at least one lowercase, one uppercase, one number and one symbol. Allowed symbols: ! @ # ~ $ % ^ & * ( ) + | _"> <br> <br>
                
                <label for="password-confirm">Confirmer le mot de passe</label>
                <input type="password" id="password-confirm" placeholder="Mot de passe" maxlength="64" required>

                <div class="error" id="error-message"></div>
                
                <button type="submit">Inscription</button>
            </form>
        `;
    }

    async init() {
        const signUpForm = document.getElementById("sign-up-form")

        signUpForm.addEventListener("submit", function () {
            const password = document.getElementById("password")
            const passwordConfirm = document.getElementById("password-confirm")

            if (password.value != passwordConfirm.value) {
                drawError("Passwords Don't Match")
            } else {
                let input = {
                    username: document.getElementById("username").value,
                    firstName: document.getElementById("first-name").value,
                    lastName: document.getElementById("last-name").value,
                    age: parseInt(document.getElementById("age").value),
                    gender: parseInt(document.querySelector('input[name="gender"]:checked').value),
                    email: document.getElementById("email").value,
                    password: password.value,
                }

                signUp(input)
            }
        })
    }
}