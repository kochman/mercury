var app = new Vue({
    el: '#app',

    data: {
       messages: [],
       contacts: [],
    },
    mounted(){
        this.fetchContacts().then(this.getMessages())
    },
    methods: {
        contactNameByID(id){
            for(let i = 0; i < this.contacts.length; i++ ){
                if (this.contacts[i].PublicKey == id){
                    return this.contacts[i].Name
                }
            }
            return "Unknown"
        },
        fetchContacts(){
            return fetch("/api/contacts/all").then((data) => data.json()).then((val) => {
                this.contacts = val;
            })
        },
        getMessages(){
            fetch("/api/messages").then((data) => data.json()).then((val) => {
                this.messages = [];
                for(let i = 0; i < val.length; i ++){
                    try{
                        val[i].Contents = JSON.parse(val[i].Contents)
                        // for(let i = 0; i < this.contacts.length; i ++){
                        //     if(contacts[i].ID == val[i].Contents.From){
                        //         val[i].Contents.From = contacts[i].Name;
                        //     }
                        // }
                        this.messages.push(val[i])
                    }catch(all){

                    }
                }
            })
        }
    },
    watch: {

    }
});