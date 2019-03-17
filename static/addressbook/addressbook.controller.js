var app = new Vue({
    el: '#app',

    data: {
        modalActive: false,
        createModal: false,
        contacts: [],
        targetContact: {},
        newContact: {
            ID: -1,
            Name: "",
            PublicKey: "",
        },
        myPubKey: "",
        myKeyModal: false,
    },
    mounted(){
        this.fetchContacts();
        fetch("/api/self").then((data)=> data.text()).then((val) => {
            this.myPubKey = val;
        })
    },
    methods: {
        fetchContacts(){
            fetch("/api/contacts/all").then((data) => data.json()).then((val) => {
                this.contacts = val;
            });
        },
        sendToUser(contact){
            window.location = "/?user=" + String(contact.ID);
        },  
        createNewContact(){
            fetch("/api/contacts/create", {
                method: "POST",
                body: JSON.stringify(this.newContact),
                
            }).then((ret) => {
                this.fetchContacts()
            })
        }
    }
});