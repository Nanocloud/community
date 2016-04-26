import Ember from 'ember';

export default Ember.Service.extend({

  downloadFile(accessToken, filename) {

    Ember.$.ajax({
      type: "GET",
      headers: { Authorization : "Bearer " + this.get('session.access_token')},
      url: "/api/files/token",
      data: { filename: "./" + filename}
    })
    .then((response) => {
      let url = "/api/files?filename=" + encodeURIComponent("./" + filename) + "&token=" + encodeURIComponent(response.token); 
      window.open(url);
    });
  }
});
