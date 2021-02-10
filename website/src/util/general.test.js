import {
  setNamedProperties,
  redirectTo,
  reduceArrayToObject,
} from './general';

jest.mock('./environment.js');

describe('setNamedProperties', () => {
  describe('success cases', () => {
    it('should copy the properties from one object to another', () => {
      const target = {};
      const source = {
        someProperty: true,
        someOtherProperty: 'ok',
        aDifferentProperty: 5,
      };
      const propertiesToSet = ['someProperty', 'someOtherProperty'];

      setNamedProperties(target, source, propertiesToSet);
      expect(target).toEqual({
        someProperty: source.someProperty,
        someOtherProperty: source.someOtherProperty,
      });
    });

    it('should camelize object property names during the copy', () => {
      const target = {};
      const source = {
        some_property: true,
        some_other_property: 'ok',
        a_different_property: 5,
      };
      const propertiesToSet = ['someProperty', 'someOtherProperty'];

      setNamedProperties(target, source, propertiesToSet);
      expect(target).toEqual({
        someProperty: source.someProperty,
        someOtherProperty: source.someOtherProperty,
      });
    });

    it('should not camelize object property names if passed false', () => {
      const target = {};
      const source = {
        some_property: true,
        some_other_property: 'ok',
        a_different_property: 5,
      };
      const propertiesToSet = ['someProperty', 'someOtherProperty'];
      const camelize = false;

      setNamedProperties(target, source, propertiesToSet, camelize);
      expect(target).toEqual({
        some_property: source.someProperty,
        some_other_property: source.someOtherProperty,
      });
    });
  });

  describe('error and null cases', () => {
    it('should do nothing if source is not an object', () => {
      const target = {};
      const source = 'not an object';
      const propertiesToSet = ['someProperty', 'someOtherProperty'];

      setNamedProperties(target, source, propertiesToSet);
      expect(target).toEqual(target);
    });

    it('should do nothing if target is not an object', () => {
      const target = 'not an object';
      const source = {
        someProperty: true,
        someOtherProperty: 'ok',
        aDifferentProperty: 5,
      };
      const propertiesToSet = ['someProperty', 'someOtherProperty'];

      setNamedProperties(target, source, propertiesToSet);
      expect(target).toEqual(target);
    });

    it('should do nothing if propertiesToSet is null', () => {
      const target = {};
      const source = {
        someProperty: true,
        someOtherProperty: 'ok',
        aDifferentProperty: 5,
      };
      const propertiesToSet = null;

      setNamedProperties(target, source, propertiesToSet);
      expect(target).toEqual(target);
    });
  });
});

describe('redirectTo', () => {
  let originalLocationAssign = null;

  beforeEach(() => {
    originalLocationAssign = window.location.assign;
    Object.defineProperty(window, 'location', {
      writable: true,
      value: { assign: jest.fn() },
    });
  });

  afterEach(() => {
    window.location.assign = originalLocationAssign;
  });

  it('should call assign on window.location', () => {
    const href = 'some href';
    redirectTo(href);
    expect(window.location.assign).toHaveBeenCalledTimes(1);
    expect(window.location.assign).toHaveBeenCalledWith(href);
  });
});

describe('reduceArrayToObject', () => {
  it('should turn array of objects into object of objects', () => {
    const arrayItem1 = { id: 5, propA: 'ok' };
    const arrayItem2 = { id: 6, propB: 'aok' };
    const arr = [arrayItem1, arrayItem2];
    const obj = arr.reduce(reduceArrayToObject, {});

    expect(obj[arrayItem1.id]).toEqual(arrayItem1);
    expect(obj[arrayItem2.id]).toEqual(arrayItem2);
  });

  it('should throw if the objects in the array do not have id', () => {
    const arrayItem1 = { propA: 'ok' };
    const arrayItem2 = { propB: 'aok' };
    const arr = [arrayItem1, arrayItem2];
    expect(() => arr.reduce(reduceArrayToObject, {})).toThrow();
  });

  it('should throw if the objects in the array are null', () => {
    const arrayItem1 = null;
    const arrayItem2 = null;
    const arr = [arrayItem1, arrayItem2];
    expect(() => arr.reduce(reduceArrayToObject, {})).toThrow();
  });
});
