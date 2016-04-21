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
    this.route('machines', function() {
      this.route('new');
      this.route('machine', { path: '/:machine_id' });
    });
    this.route('apps', function() {
      this.route('app', { path: '/:app_id' });
    });
    this.route('files', function() {
      this.route('nowindows');
      this.route('notactivated');
    });
    this.route('histories', function() {});
  });

  this.route('login');
});

export default Router;
