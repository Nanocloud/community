import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('services', {path: '/'});
  this.route('users');

  this.route('user', function() {
    this.route('create');
  });
});

export default Router;
