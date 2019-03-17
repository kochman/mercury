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
        myName: "joey",

    },
    mounted(){
        this.fetchContacts();
        this.getMyName();
        fetch("/api/self").then((data)=> data.text()).then((val) => {
            this.myPubKey = val;
        })
    },
    methods: {
        sendMyName(){
            fetch("/myinfo", {
                method: "POST",
                body: JSON.stringify({Name: this.myName})
            }).then(this.getMyName)
        },
        getMyName(){
            fetch("/myinfo").then((val) => val.json()).then((data) => {
                this.myName = data.Name
            })
        },
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
                this.createModal = false;
            })
        }
    }
});