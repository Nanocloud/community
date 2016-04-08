import Ember from 'ember';

export default Ember.Controller.extend({

  store: Ember.inject.service('store'),
  session: Ember.inject.service('session'),
  items: null,

  loadDir: function() {
    let path = "./";
    this.get('store').query('file', { filename: path })
      .then((response) => {
        this.set('items', response);
      })
  }.on('init'),

  actions : {

    downloadFile: function(filename) {

     Ember.$.ajax({
        type: "GET",
        headers: { Authorization : "Bearer " + this.get('session.access_token') },
        url: "/api/files/token",
        data: { filename: "./" + filename}
      })
     .then((response) => {
        let url = "/api/files?filename=" + encodeURIComponent("./" + filename) + "&token=" + encodeURIComponent(response.token); 
        window.open(url);
     }, () => {
       this.toast.error("Couldn't retrieve a token to download the file");
     });
    },
  }
});
