import fetcher from "../services/Fetcher.js";
import AbstractView from "./AbstractView.js";

const genders = { 1: 'Homme', 2: 'Femme' }

const getUserByID = async (id) => {
    const path = `/api/users/${id}`
    return await fetcher.get(path);
}

const getUsersPosts = async (userID) => {
    const path = `/api/users/${userID}/posts`
    return await fetcher.get(path);
}

const getUsersRatedPosts = async (userID) => {
    const path = `/api/users/${userID}/rated-posts`
    return await fetcher.get(path);
}

const newPostElement = (post) => {
    const el = document.createElement("div")
    el.classList.add("post")

    const linkToPost = document.createElement("a")
    linkToPost.classList.add("post-link")
    linkToPost.setAttribute("href", `/post/${post.id}`)
    linkToPost.setAttribute("data-link", "")
    linkToPost.innerText = `${post.title}`

    const postDate = document.createElement("p")
    postDate.innerText = new Date(post.date).toLocaleString()

    const linkToAuthor = document.createElement("a")
    linkToAuthor.setAttribute("href", `/user/${post.author.id}`)
    linkToAuthor.setAttribute("data-link", "")
    linkToAuthor.innerText = `${post.author.firstName} ${post.author.lastName}`

    el.append(linkToPost)
    el.append(postDate)
    el.append(linkToAuthor)

    return el
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Profile");
        this.userID = params.userID;
    }

    async getHtml() {
        return `
        <h2>Mon profil</h2>
        <div id="user-profile">
                <div class="profile-info" id="avatar"></div>
                <div>
                    <div class="profile-info" id="username"></div>
                    <div class="profile-info" id="first-name"></div>
                    <div class="profile-info" id="last-name"></div>
                    <div class="profile-info" id="age"></div>
                    <div class="profile-info" id="gender"></div>
                    <div class="profile-info" id="registered"></div>
                </div>
            </div>
            <h2>Posts publiés</h2>
            <div id="users-posts"></div>
            <h2>Posts que vous avez aimé</h2>
            <div id="users-liked-posts"></div>
        `;
    }

    async init() {
        const user = await getUserByID(this.userID)
       
        document.querySelector('.profile-info#avatar').innerHTML = `<img src="http://${API_HOST_NAME}/images/${user.avatar}">`
        document.querySelector('.profile-info#username').innerText = `Pseudo : ${user.username}`
        document.querySelector('.profile-info#first-name').innerText = `Prénom : ${user.firstName}`
        document.querySelector('.profile-info#last-name').innerText = `Nom : ${user.lastName}`
        document.querySelector('.profile-info#age').innerText = `Age: ${user.age}`
        document.querySelector('.profile-info#gender').innerText = `Genre : ${genders[user.gender]}`
        document.querySelector('.profile-info#registered').innerText = `Inscrit le  ${new Date(Date.parse(user.registered)).toLocaleString()}`

        const usersPosts = await getUsersPosts(this.userID) 
        const usersRatedPosts = await getUsersRatedPosts(this.userID)|| []

        const usersPostsEl = document.getElementById('users-posts')
        if (usersPosts != null) {
            usersPosts.forEach((post) => {
                const postEl = newPostElement(post)
                usersPostsEl.append(postEl)
            })
        } else {
            usersPostsEl.innerText = 'Aucun posts'
        }


        const usersLikedPosts = usersRatedPosts.filter((post) => post.userRate == 1)
        const usersLikedPostsEl = document.getElementById('users-liked-posts')
        if (usersLikedPosts.length > 0 ) {
            usersLikedPosts.forEach((post) => {
                const postEl = newPostElement(post)
                usersLikedPostsEl.append(postEl)
            })
        } else {
            usersLikedPostsEl.innerText = 'Aucun posts'
        }

        const usersDisLikedPosts = usersRatedPosts.filter((post) => post.userRate == 2)
        const usersDislikedPostsEl = document.getElementById('users-disliked-posts')
        if (usersDisLikedPosts.length > 0) {
            usersDisLikedPosts.forEach((post) => {
                const postEl = newPostElement(post)
                usersDislikedPostsEl.append(postEl)
            })
        } else {
            usersDislikedPostsEl.innerText = 'Aucun posts'
        }
    }
}