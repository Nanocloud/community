import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('protected', {path: '/'}, function() {
    this.route('users', function() {
      this.route('user', { path: '/:user_id' });
      this.route('new');
    });
    this.route('machines');
    this.route('applications');
  });

  this.route('login');
});

export default Router;
