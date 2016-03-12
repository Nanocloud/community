import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('protected', {path: '/'}, function() {
    this.route('services');
    this.route('users');

    this.route('user', function() {
      this.route('create');
    });

  });

  this.route('login');
});

export default Router;
