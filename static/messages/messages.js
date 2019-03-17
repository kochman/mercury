var app = new Vue({
    el: '#app',

    data: {
       messages: [],
       contacts: [],
    },
    mounted(){
        setInterval(this.getMessages, 2000)
        this.fetchContacts().then(this.getMessages())
    },
    computed: {
        msgArr: function() {
            this.messages.sort( ( a, b) => {
                return new Date(b.Sent) - new Date(a.Sent);
            });
            return this.messages;
        }
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
            console.log("here")
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