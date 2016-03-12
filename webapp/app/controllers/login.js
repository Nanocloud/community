import Ember from 'ember';

export default Ember.Controller.extend({
  session: Ember.inject.service('session'),

  actions: {
    authenticate: function() {
      let { identification, password } = this.getProperties('identification', 'password');

      this.get('session').authenticate('authenticator:oauth2', identification, password)
        .then((user) => {
          this.transitionToRoute('/');
        })
        .catch((reason) => {
          this.set('errorMessage', reason.error || reason);
      });
    }
  }
});
