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
        fetch("/self").then((data)=> data.text()).then((val) => {
            this.myPubKey = val;
        })
    },
    methods: {
        fetchContacts(){
            fetch("/contacts/all").then((data) => data.json()).then((val) => {
                this.contacts = val;
            });
        },
        createNewContact(){
            fetch("/contacts/create", {
                method: "POST",
                body: JSON.stringify(this.newContact),
                
            }).then((ret) => {
                this.fetchContacts()
            })
        }
    }
});
