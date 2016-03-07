import Ember from 'ember';
import DS from 'ember-data';

var underscore = Ember.String.underscore;

// Allow to have underscore in JSONAPI keys
export default DS.JSONAPISerializer.extend({
  keyForAttribute: function(attr) {
    return underscore(attr);
  },

  keyForRelationship: function(rawKey) {
    return underscore(rawKey);
  }
});
