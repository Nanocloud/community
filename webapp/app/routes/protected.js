import Ember from 'ember';

export default Ember.Route.extend({
  session: Ember.inject.service('session'),

  beforeModel() {
    if (this.get('session.isAuthenticated') === false) {
      this.transitionTo('login');
    }
  },

  model: function() {
    this.set('session.access_token', this.get('session.data.authenticated.access_token'));

    return new Ember.RSVP.Promise((resolve, reject) => {
      const options = {
        url: 'api/me',
        headers: {
          Authorization: 'Bearer ' + this.get('session.access_token')
        },
        type:        'GET',
        dataType:    'json'
      };

      Ember.$.ajax(options).then((user) => {
        this.set('session.user', user.data);
        resolve(user);
      }, (issue) => {
        this._invalidateSession();
        reject(issue);
      });
    });
  },

  _invalidateSession: function() {
      this.get('session').invalidate();
      this.transitionTo('login');
  },

  actions: {
    invalidateSession: function() {
      this._invalidateSession();
    }
  }
});
