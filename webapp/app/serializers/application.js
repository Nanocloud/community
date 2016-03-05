import Ember from 'ember';
import DS from 'ember-data';

var underscore = Ember.String.underscore;

// Allow to have understand in JSONAPI keys
export default DS.JSONAPISerializer.extend({
  keyForAttribute: function(attr) {
    return underscore(attr);
  },

  keyForRelationship: function(rawKey) {
    return underscore(rawKey);
  }
});
