import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
    this.route('canva', {path: '/canva/:token/:connectionName'});
});

export default Router;
