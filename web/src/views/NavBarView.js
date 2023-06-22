import AbstractView from "./AbstractView.js";
import Utils from "../services/Utils.js"
import Ws from "../services/Ws.js";
import intervals from "../services/Intervals.js";


export default class extends AbstractView {
    constructor(params, user) {
        super(params);
        this.user = user
    }

    async getHtml() {
        const authorized = Boolean(this.user.id)

        return `
            <img class="logoforum" src="src/img/lettre-f.png" width="50">
            <div class="nav__links">
            <a href="/" class="nav__link" id="home-button" data-link>Accueil</a>
            `
            +
            (authorized ?
                `
                <a href="/new-post" class="nav__link" id="new-post-button" data-link>Créer un post</a>
                <a href="/chats" class="nav__link" id="chats-button" data-link>Messagerie</a>
                <a href="/user/${this.user.id}" class="nav__link" id="chats-button" data-link>Mon profil</a>
                <a href="/" class="nav__link" id="sign-out-button" data-link>Déconnexion</a>
                `
                :
                `
                <a href="/sign-up" class="nav__link" id="sign-in-button" data-link>Inscription</a>
                <a href="/sign-in" class="nav__link" id="sign-up-button" data-link>Connexion</a>
                </div>
                `
            )
    }

    async init() {
        const signOutButton = document.getElementById("sign-out-button")
        if (signOutButton) {
            signOutButton.addEventListener('click', () => {
                Utils.logOut()
                Ws.disconnect()
                intervals.clearAll()
            })
        }
    }
}