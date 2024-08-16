(() => {
  // node_modules/.pnpm/alpinejs@3.14.1/node_modules/alpinejs/dist/module.esm.js
  var flushPending = false;
  var flushing = false;
  var queue = [];
  var lastFlushedIndex = -1;
  function scheduler(callback) {
    queueJob(callback);
  }
  function queueJob(job) {
    if (!queue.includes(job))
      queue.push(job);
    queueFlush();
  }
  function dequeueJob(job) {
    let index = queue.indexOf(job);
    if (index !== -1 && index > lastFlushedIndex)
      queue.splice(index, 1);
  }
  function queueFlush() {
    if (!flushing && !flushPending) {
      flushPending = true;
      queueMicrotask(flushJobs);
    }
  }
  function flushJobs() {
    flushPending = false;
    flushing = true;
    for (let i = 0; i < queue.length; i++) {
      queue[i]();
      lastFlushedIndex = i;
    }
    queue.length = 0;
    lastFlushedIndex = -1;
    flushing = false;
  }
  var reactive;
  var effect;
  var release;
  var raw;
  var shouldSchedule = true;
  function disableEffectScheduling(callback) {
    shouldSchedule = false;
    callback();
    shouldSchedule = true;
  }
  function setReactivityEngine(engine) {
    reactive = engine.reactive;
    release = engine.release;
    effect = (callback) => engine.effect(callback, { scheduler: (task) => {
      if (shouldSchedule) {
        scheduler(task);
      } else {
        task();
      }
    } });
    raw = engine.raw;
  }
  function overrideEffect(override) {
    effect = override;
  }
  function elementBoundEffect(el) {
    let cleanup2 = () => {
    };
    let wrappedEffect = (callback) => {
      let effectReference = effect(callback);
      if (!el._x_effects) {
        el._x_effects = /* @__PURE__ */ new Set();
        el._x_runEffects = () => {
          el._x_effects.forEach((i) => i());
        };
      }
      el._x_effects.add(effectReference);
      cleanup2 = () => {
        if (effectReference === void 0)
          return;
        el._x_effects.delete(effectReference);
        release(effectReference);
      };
      return effectReference;
    };
    return [wrappedEffect, () => {
      cleanup2();
    }];
  }
  function watch(getter, callback) {
    let firstTime = true;
    let oldValue;
    let effectReference = effect(() => {
      let value = getter();
      JSON.stringify(value);
      if (!firstTime) {
        queueMicrotask(() => {
          callback(value, oldValue);
          oldValue = value;
        });
      } else {
        oldValue = value;
      }
      firstTime = false;
    });
    return () => release(effectReference);
  }
  var onAttributeAddeds = [];
  var onElRemoveds = [];
  var onElAddeds = [];
  function onElAdded(callback) {
    onElAddeds.push(callback);
  }
  function onElRemoved(el, callback) {
    if (typeof callback === "function") {
      if (!el._x_cleanups)
        el._x_cleanups = [];
      el._x_cleanups.push(callback);
    } else {
      callback = el;
      onElRemoveds.push(callback);
    }
  }
  function onAttributesAdded(callback) {
    onAttributeAddeds.push(callback);
  }
  function onAttributeRemoved(el, name, callback) {
    if (!el._x_attributeCleanups)
      el._x_attributeCleanups = {};
    if (!el._x_attributeCleanups[name])
      el._x_attributeCleanups[name] = [];
    el._x_attributeCleanups[name].push(callback);
  }
  function cleanupAttributes(el, names) {
    if (!el._x_attributeCleanups)
      return;
    Object.entries(el._x_attributeCleanups).forEach(([name, value]) => {
      if (names === void 0 || names.includes(name)) {
        value.forEach((i) => i());
        delete el._x_attributeCleanups[name];
      }
    });
  }
  function cleanupElement(el) {
    if (el._x_cleanups) {
      while (el._x_cleanups.length)
        el._x_cleanups.pop()();
    }
  }
  var observer = new MutationObserver(onMutate);
  var currentlyObserving = false;
  function startObservingMutations() {
    observer.observe(document, { subtree: true, childList: true, attributes: true, attributeOldValue: true });
    currentlyObserving = true;
  }
  function stopObservingMutations() {
    flushObserver();
    observer.disconnect();
    currentlyObserving = false;
  }
  var queuedMutations = [];
  function flushObserver() {
    let records = observer.takeRecords();
    queuedMutations.push(() => records.length > 0 && onMutate(records));
    let queueLengthWhenTriggered = queuedMutations.length;
    queueMicrotask(() => {
      if (queuedMutations.length === queueLengthWhenTriggered) {
        while (queuedMutations.length > 0)
          queuedMutations.shift()();
      }
    });
  }
  function mutateDom(callback) {
    if (!currentlyObserving)
      return callback();
    stopObservingMutations();
    let result = callback();
    startObservingMutations();
    return result;
  }
  var isCollecting = false;
  var deferredMutations = [];
  function deferMutations() {
    isCollecting = true;
  }
  function flushAndStopDeferringMutations() {
    isCollecting = false;
    onMutate(deferredMutations);
    deferredMutations = [];
  }
  function onMutate(mutations) {
    if (isCollecting) {
      deferredMutations = deferredMutations.concat(mutations);
      return;
    }
    let addedNodes = /* @__PURE__ */ new Set();
    let removedNodes = /* @__PURE__ */ new Set();
    let addedAttributes = /* @__PURE__ */ new Map();
    let removedAttributes = /* @__PURE__ */ new Map();
    for (let i = 0; i < mutations.length; i++) {
      if (mutations[i].target._x_ignoreMutationObserver)
        continue;
      if (mutations[i].type === "childList") {
        mutations[i].addedNodes.forEach((node) => node.nodeType === 1 && addedNodes.add(node));
        mutations[i].removedNodes.forEach((node) => node.nodeType === 1 && removedNodes.add(node));
      }
      if (mutations[i].type === "attributes") {
        let el = mutations[i].target;
        let name = mutations[i].attributeName;
        let oldValue = mutations[i].oldValue;
        let add2 = () => {
          if (!addedAttributes.has(el))
            addedAttributes.set(el, []);
          addedAttributes.get(el).push({ name, value: el.getAttribute(name) });
        };
        let remove = () => {
          if (!removedAttributes.has(el))
            removedAttributes.set(el, []);
          removedAttributes.get(el).push(name);
        };
        if (el.hasAttribute(name) && oldValue === null) {
          add2();
        } else if (el.hasAttribute(name)) {
          remove();
          add2();
        } else {
          remove();
        }
      }
    }
    removedAttributes.forEach((attrs, el) => {
      cleanupAttributes(el, attrs);
    });
    addedAttributes.forEach((attrs, el) => {
      onAttributeAddeds.forEach((i) => i(el, attrs));
    });
    for (let node of removedNodes) {
      if (addedNodes.has(node))
        continue;
      onElRemoveds.forEach((i) => i(node));
    }
    addedNodes.forEach((node) => {
      node._x_ignoreSelf = true;
      node._x_ignore = true;
    });
    for (let node of addedNodes) {
      if (removedNodes.has(node))
        continue;
      if (!node.isConnected)
        continue;
      delete node._x_ignoreSelf;
      delete node._x_ignore;
      onElAddeds.forEach((i) => i(node));
      node._x_ignore = true;
      node._x_ignoreSelf = true;
    }
    addedNodes.forEach((node) => {
      delete node._x_ignoreSelf;
      delete node._x_ignore;
    });
    addedNodes = null;
    removedNodes = null;
    addedAttributes = null;
    removedAttributes = null;
  }
  function scope(node) {
    return mergeProxies(closestDataStack(node));
  }
  function addScopeToNode(node, data2, referenceNode) {
    node._x_dataStack = [data2, ...closestDataStack(referenceNode || node)];
    return () => {
      node._x_dataStack = node._x_dataStack.filter((i) => i !== data2);
    };
  }
  function closestDataStack(node) {
    if (node._x_dataStack)
      return node._x_dataStack;
    if (typeof ShadowRoot === "function" && node instanceof ShadowRoot) {
      return closestDataStack(node.host);
    }
    if (!node.parentNode) {
      return [];
    }
    return closestDataStack(node.parentNode);
  }
  function mergeProxies(objects) {
    return new Proxy({ objects }, mergeProxyTrap);
  }
  var mergeProxyTrap = {
    ownKeys({ objects }) {
      return Array.from(
        new Set(objects.flatMap((i) => Object.keys(i)))
      );
    },
    has({ objects }, name) {
      if (name == Symbol.unscopables)
        return false;
      return objects.some(
        (obj) => Object.prototype.hasOwnProperty.call(obj, name) || Reflect.has(obj, name)
      );
    },
    get({ objects }, name, thisProxy) {
      if (name == "toJSON")
        return collapseProxies;
      return Reflect.get(
        objects.find(
          (obj) => Reflect.has(obj, name)
        ) || {},
        name,
        thisProxy
      );
    },
    set({ objects }, name, value, thisProxy) {
      const target = objects.find(
        (obj) => Object.prototype.hasOwnProperty.call(obj, name)
      ) || objects[objects.length - 1];
      const descriptor = Object.getOwnPropertyDescriptor(target, name);
      if (descriptor?.set && descriptor?.get)
        return descriptor.set.call(thisProxy, value) || true;
      return Reflect.set(target, name, value);
    }
  };
  function collapseProxies() {
    let keys = Reflect.ownKeys(this);
    return keys.reduce((acc, key) => {
      acc[key] = Reflect.get(this, key);
      return acc;
    }, {});
  }
  function initInterceptors(data2) {
    let isObject2 = (val) => typeof val === "object" && !Array.isArray(val) && val !== null;
    let recurse = (obj, basePath = "") => {
      Object.entries(Object.getOwnPropertyDescriptors(obj)).forEach(([key, { value, enumerable }]) => {
        if (enumerable === false || value === void 0)
          return;
        if (typeof value === "object" && value !== null && value.__v_skip)
          return;
        let path = basePath === "" ? key : `${basePath}.${key}`;
        if (typeof value === "object" && value !== null && value._x_interceptor) {
          obj[key] = value.initialize(data2, path, key);
        } else {
          if (isObject2(value) && value !== obj && !(value instanceof Element)) {
            recurse(value, path);
          }
        }
      });
    };
    return recurse(data2);
  }
  function interceptor(callback, mutateObj = () => {
  }) {
    let obj = {
      initialValue: void 0,
      _x_interceptor: true,
      initialize(data2, path, key) {
        return callback(this.initialValue, () => get(data2, path), (value) => set(data2, path, value), path, key);
      }
    };
    mutateObj(obj);
    return (initialValue) => {
      if (typeof initialValue === "object" && initialValue !== null && initialValue._x_interceptor) {
        let initialize = obj.initialize.bind(obj);
        obj.initialize = (data2, path, key) => {
          let innerValue = initialValue.initialize(data2, path, key);
          obj.initialValue = innerValue;
          return initialize(data2, path, key);
        };
      } else {
        obj.initialValue = initialValue;
      }
      return obj;
    };
  }
  function get(obj, path) {
    return path.split(".").reduce((carry, segment) => carry[segment], obj);
  }
  function set(obj, path, value) {
    if (typeof path === "string")
      path = path.split(".");
    if (path.length === 1)
      obj[path[0]] = value;
    else if (path.length === 0)
      throw error;
    else {
      if (obj[path[0]])
        return set(obj[path[0]], path.slice(1), value);
      else {
        obj[path[0]] = {};
        return set(obj[path[0]], path.slice(1), value);
      }
    }
  }
  var magics = {};
  function magic(name, callback) {
    magics[name] = callback;
  }
  function injectMagics(obj, el) {
    Object.entries(magics).forEach(([name, callback]) => {
      let memoizedUtilities = null;
      function getUtilities() {
        if (memoizedUtilities) {
          return memoizedUtilities;
        } else {
          let [utilities, cleanup2] = getElementBoundUtilities(el);
          memoizedUtilities = { interceptor, ...utilities };
          onElRemoved(el, cleanup2);
          return memoizedUtilities;
        }
      }
      Object.defineProperty(obj, `$${name}`, {
        get() {
          return callback(el, getUtilities());
        },
        enumerable: false
      });
    });
    return obj;
  }
  function tryCatch(el, expression, callback, ...args) {
    try {
      return callback(...args);
    } catch (e) {
      handleError(e, el, expression);
    }
  }
  function handleError(error2, el, expression = void 0) {
    error2 = Object.assign(
      error2 ?? { message: "No error message given." },
      { el, expression }
    );
    console.warn(`Alpine Expression Error: ${error2.message}

${expression ? 'Expression: "' + expression + '"\n\n' : ""}`, el);
    setTimeout(() => {
      throw error2;
    }, 0);
  }
  var shouldAutoEvaluateFunctions = true;
  function dontAutoEvaluateFunctions(callback) {
    let cache = shouldAutoEvaluateFunctions;
    shouldAutoEvaluateFunctions = false;
    let result = callback();
    shouldAutoEvaluateFunctions = cache;
    return result;
  }
  function evaluate(el, expression, extras = {}) {
    let result;
    evaluateLater(el, expression)((value) => result = value, extras);
    return result;
  }
  function evaluateLater(...args) {
    return theEvaluatorFunction(...args);
  }
  var theEvaluatorFunction = normalEvaluator;
  function setEvaluator(newEvaluator) {
    theEvaluatorFunction = newEvaluator;
  }
  function normalEvaluator(el, expression) {
    let overriddenMagics = {};
    injectMagics(overriddenMagics, el);
    let dataStack = [overriddenMagics, ...closestDataStack(el)];
    let evaluator = typeof expression === "function" ? generateEvaluatorFromFunction(dataStack, expression) : generateEvaluatorFromString(dataStack, expression, el);
    return tryCatch.bind(null, el, expression, evaluator);
  }
  function generateEvaluatorFromFunction(dataStack, func) {
    return (receiver = () => {
    }, { scope: scope2 = {}, params = [] } = {}) => {
      let result = func.apply(mergeProxies([scope2, ...dataStack]), params);
      runIfTypeOfFunction(receiver, result);
    };
  }
  var evaluatorMemo = {};
  function generateFunctionFromString(expression, el) {
    if (evaluatorMemo[expression]) {
      return evaluatorMemo[expression];
    }
    let AsyncFunction = Object.getPrototypeOf(async function() {
    }).constructor;
    let rightSideSafeExpression = /^[\n\s]*if.*\(.*\)/.test(expression.trim()) || /^(let|const)\s/.test(expression.trim()) ? `(async()=>{ ${expression} })()` : expression;
    const safeAsyncFunction = () => {
      try {
        let func2 = new AsyncFunction(
          ["__self", "scope"],
          `with (scope) { __self.result = ${rightSideSafeExpression} }; __self.finished = true; return __self.result;`
        );
        Object.defineProperty(func2, "name", {
          value: `[Alpine] ${expression}`
        });
        return func2;
      } catch (error2) {
        handleError(error2, el, expression);
        return Promise.resolve();
      }
    };
    let func = safeAsyncFunction();
    evaluatorMemo[expression] = func;
    return func;
  }
  function generateEvaluatorFromString(dataStack, expression, el) {
    let func = generateFunctionFromString(expression, el);
    return (receiver = () => {
    }, { scope: scope2 = {}, params = [] } = {}) => {
      func.result = void 0;
      func.finished = false;
      let completeScope = mergeProxies([scope2, ...dataStack]);
      if (typeof func === "function") {
        let promise = func(func, completeScope).catch((error2) => handleError(error2, el, expression));
        if (func.finished) {
          runIfTypeOfFunction(receiver, func.result, completeScope, params, el);
          func.result = void 0;
        } else {
          promise.then((result) => {
            runIfTypeOfFunction(receiver, result, completeScope, params, el);
          }).catch((error2) => handleError(error2, el, expression)).finally(() => func.result = void 0);
        }
      }
    };
  }
  function runIfTypeOfFunction(receiver, value, scope2, params, el) {
    if (shouldAutoEvaluateFunctions && typeof value === "function") {
      let result = value.apply(scope2, params);
      if (result instanceof Promise) {
        result.then((i) => runIfTypeOfFunction(receiver, i, scope2, params)).catch((error2) => handleError(error2, el, value));
      } else {
        receiver(result);
      }
    } else if (typeof value === "object" && value instanceof Promise) {
      value.then((i) => receiver(i));
    } else {
      receiver(value);
    }
  }
  var prefixAsString = "x-";
  function prefix(subject = "") {
    return prefixAsString + subject;
  }
  function setPrefix(newPrefix) {
    prefixAsString = newPrefix;
  }
  var directiveHandlers = {};
  function directive(name, callback) {
    directiveHandlers[name] = callback;
    return {
      before(directive2) {
        if (!directiveHandlers[directive2]) {
          console.warn(String.raw`Cannot find directive \`${directive2}\`. \`${name}\` will use the default order of execution`);
          return;
        }
        const pos = directiveOrder.indexOf(directive2);
        directiveOrder.splice(pos >= 0 ? pos : directiveOrder.indexOf("DEFAULT"), 0, name);
      }
    };
  }
  function directiveExists(name) {
    return Object.keys(directiveHandlers).includes(name);
  }
  function directives(el, attributes, originalAttributeOverride) {
    attributes = Array.from(attributes);
    if (el._x_virtualDirectives) {
      let vAttributes = Object.entries(el._x_virtualDirectives).map(([name, value]) => ({ name, value }));
      let staticAttributes = attributesOnly(vAttributes);
      vAttributes = vAttributes.map((attribute) => {
        if (staticAttributes.find((attr) => attr.name === attribute.name)) {
          return {
            name: `x-bind:${attribute.name}`,
            value: `"${attribute.value}"`
          };
        }
        return attribute;
      });
      attributes = attributes.concat(vAttributes);
    }
    let transformedAttributeMap = {};
    let directives2 = attributes.map(toTransformedAttributes((newName, oldName) => transformedAttributeMap[newName] = oldName)).filter(outNonAlpineAttributes).map(toParsedDirectives(transformedAttributeMap, originalAttributeOverride)).sort(byPriority);
    return directives2.map((directive2) => {
      return getDirectiveHandler(el, directive2);
    });
  }
  function attributesOnly(attributes) {
    return Array.from(attributes).map(toTransformedAttributes()).filter((attr) => !outNonAlpineAttributes(attr));
  }
  var isDeferringHandlers = false;
  var directiveHandlerStacks = /* @__PURE__ */ new Map();
  var currentHandlerStackKey = Symbol();
  function deferHandlingDirectives(callback) {
    isDeferringHandlers = true;
    let key = Symbol();
    currentHandlerStackKey = key;
    directiveHandlerStacks.set(key, []);
    let flushHandlers = () => {
      while (directiveHandlerStacks.get(key).length)
        directiveHandlerStacks.get(key).shift()();
      directiveHandlerStacks.delete(key);
    };
    let stopDeferring = () => {
      isDeferringHandlers = false;
      flushHandlers();
    };
    callback(flushHandlers);
    stopDeferring();
  }
  function getElementBoundUtilities(el) {
    let cleanups = [];
    let cleanup2 = (callback) => cleanups.push(callback);
    let [effect3, cleanupEffect] = elementBoundEffect(el);
    cleanups.push(cleanupEffect);
    let utilities = {
      Alpine: alpine_default,
      effect: effect3,
      cleanup: cleanup2,
      evaluateLater: evaluateLater.bind(evaluateLater, el),
      evaluate: evaluate.bind(evaluate, el)
    };
    let doCleanup = () => cleanups.forEach((i) => i());
    return [utilities, doCleanup];
  }
  function getDirectiveHandler(el, directive2) {
    let noop = () => {
    };
    let handler4 = directiveHandlers[directive2.type] || noop;
    let [utilities, cleanup2] = getElementBoundUtilities(el);
    onAttributeRemoved(el, directive2.original, cleanup2);
    let fullHandler = () => {
      if (el._x_ignore || el._x_ignoreSelf)
        return;
      handler4.inline && handler4.inline(el, directive2, utilities);
      handler4 = handler4.bind(handler4, el, directive2, utilities);
      isDeferringHandlers ? directiveHandlerStacks.get(currentHandlerStackKey).push(handler4) : handler4();
    };
    fullHandler.runCleanups = cleanup2;
    return fullHandler;
  }
  var startingWith = (subject, replacement) => ({ name, value }) => {
    if (name.startsWith(subject))
      name = name.replace(subject, replacement);
    return { name, value };
  };
  var into = (i) => i;
  function toTransformedAttributes(callback = () => {
  }) {
    return ({ name, value }) => {
      let { name: newName, value: newValue } = attributeTransformers.reduce((carry, transform) => {
        return transform(carry);
      }, { name, value });
      if (newName !== name)
        callback(newName, name);
      return { name: newName, value: newValue };
    };
  }
  var attributeTransformers = [];
  function mapAttributes(callback) {
    attributeTransformers.push(callback);
  }
  function outNonAlpineAttributes({ name }) {
    return alpineAttributeRegex().test(name);
  }
  var alpineAttributeRegex = () => new RegExp(`^${prefixAsString}([^:^.]+)\\b`);
  function toParsedDirectives(transformedAttributeMap, originalAttributeOverride) {
    return ({ name, value }) => {
      let typeMatch = name.match(alpineAttributeRegex());
      let valueMatch = name.match(/:([a-zA-Z0-9\-_:]+)/);
      let modifiers = name.match(/\.[^.\]]+(?=[^\]]*$)/g) || [];
      let original = originalAttributeOverride || transformedAttributeMap[name] || name;
      return {
        type: typeMatch ? typeMatch[1] : null,
        value: valueMatch ? valueMatch[1] : null,
        modifiers: modifiers.map((i) => i.replace(".", "")),
        expression: value,
        original
      };
    };
  }
  var DEFAULT = "DEFAULT";
  var directiveOrder = [
    "ignore",
    "ref",
    "data",
    "id",
    "anchor",
    "bind",
    "init",
    "for",
    "model",
    "modelable",
    "transition",
    "show",
    "if",
    DEFAULT,
    "teleport"
  ];
  function byPriority(a, b) {
    let typeA = directiveOrder.indexOf(a.type) === -1 ? DEFAULT : a.type;
    let typeB = directiveOrder.indexOf(b.type) === -1 ? DEFAULT : b.type;
    return directiveOrder.indexOf(typeA) - directiveOrder.indexOf(typeB);
  }
  function dispatch(el, name, detail = {}) {
    el.dispatchEvent(
      new CustomEvent(name, {
        detail,
        bubbles: true,
        // Allows events to pass the shadow DOM barrier.
        composed: true,
        cancelable: true
      })
    );
  }
  function walk(el, callback) {
    if (typeof ShadowRoot === "function" && el instanceof ShadowRoot) {
      Array.from(el.children).forEach((el2) => walk(el2, callback));
      return;
    }
    let skip = false;
    callback(el, () => skip = true);
    if (skip)
      return;
    let node = el.firstElementChild;
    while (node) {
      walk(node, callback, false);
      node = node.nextElementSibling;
    }
  }
  function warn(message, ...args) {
    console.warn(`Alpine Warning: ${message}`, ...args);
  }
  var started = false;
  function start() {
    if (started)
      warn("Alpine has already been initialized on this page. Calling Alpine.start() more than once can cause problems.");
    started = true;
    if (!document.body)
      warn("Unable to initialize. Trying to load Alpine before `<body>` is available. Did you forget to add `defer` in Alpine's `<script>` tag?");
    dispatch(document, "alpine:init");
    dispatch(document, "alpine:initializing");
    startObservingMutations();
    onElAdded((el) => initTree(el, walk));
    onElRemoved((el) => destroyTree(el));
    onAttributesAdded((el, attrs) => {
      directives(el, attrs).forEach((handle) => handle());
    });
    let outNestedComponents = (el) => !closestRoot(el.parentElement, true);
    Array.from(document.querySelectorAll(allSelectors().join(","))).filter(outNestedComponents).forEach((el) => {
      initTree(el);
    });
    dispatch(document, "alpine:initialized");
    setTimeout(() => {
      warnAboutMissingPlugins();
    });
  }
  var rootSelectorCallbacks = [];
  var initSelectorCallbacks = [];
  function rootSelectors() {
    return rootSelectorCallbacks.map((fn) => fn());
  }
  function allSelectors() {
    return rootSelectorCallbacks.concat(initSelectorCallbacks).map((fn) => fn());
  }
  function addRootSelector(selectorCallback) {
    rootSelectorCallbacks.push(selectorCallback);
  }
  function addInitSelector(selectorCallback) {
    initSelectorCallbacks.push(selectorCallback);
  }
  function closestRoot(el, includeInitSelectors = false) {
    return findClosest(el, (element) => {
      const selectors = includeInitSelectors ? allSelectors() : rootSelectors();
      if (selectors.some((selector) => element.matches(selector)))
        return true;
    });
  }
  function findClosest(el, callback) {
    if (!el)
      return;
    if (callback(el))
      return el;
    if (el._x_teleportBack)
      el = el._x_teleportBack;
    if (!el.parentElement)
      return;
    return findClosest(el.parentElement, callback);
  }
  function isRoot(el) {
    return rootSelectors().some((selector) => el.matches(selector));
  }
  var initInterceptors2 = [];
  function interceptInit(callback) {
    initInterceptors2.push(callback);
  }
  function initTree(el, walker = walk, intercept = () => {
  }) {
    deferHandlingDirectives(() => {
      walker(el, (el2, skip) => {
        intercept(el2, skip);
        initInterceptors2.forEach((i) => i(el2, skip));
        directives(el2, el2.attributes).forEach((handle) => handle());
        el2._x_ignore && skip();
      });
    });
  }
  function destroyTree(root, walker = walk) {
    walker(root, (el) => {
      cleanupAttributes(el);
      cleanupElement(el);
    });
  }
  function warnAboutMissingPlugins() {
    let pluginDirectives = [
      ["ui", "dialog", ["[x-dialog], [x-popover]"]],
      ["anchor", "anchor", ["[x-anchor]"]],
      ["sort", "sort", ["[x-sort]"]]
    ];
    pluginDirectives.forEach(([plugin2, directive2, selectors]) => {
      if (directiveExists(directive2))
        return;
      selectors.some((selector) => {
        if (document.querySelector(selector)) {
          warn(`found "${selector}", but missing ${plugin2} plugin`);
          return true;
        }
      });
    });
  }
  var tickStack = [];
  var isHolding = false;
  function nextTick(callback = () => {
  }) {
    queueMicrotask(() => {
      isHolding || setTimeout(() => {
        releaseNextTicks();
      });
    });
    return new Promise((res) => {
      tickStack.push(() => {
        callback();
        res();
      });
    });
  }
  function releaseNextTicks() {
    isHolding = false;
    while (tickStack.length)
      tickStack.shift()();
  }
  function holdNextTicks() {
    isHolding = true;
  }
  function setClasses(el, value) {
    if (Array.isArray(value)) {
      return setClassesFromString(el, value.join(" "));
    } else if (typeof value === "object" && value !== null) {
      return setClassesFromObject(el, value);
    } else if (typeof value === "function") {
      return setClasses(el, value());
    }
    return setClassesFromString(el, value);
  }
  function setClassesFromString(el, classString) {
    let split = (classString2) => classString2.split(" ").filter(Boolean);
    let missingClasses = (classString2) => classString2.split(" ").filter((i) => !el.classList.contains(i)).filter(Boolean);
    let addClassesAndReturnUndo = (classes) => {
      el.classList.add(...classes);
      return () => {
        el.classList.remove(...classes);
      };
    };
    classString = classString === true ? classString = "" : classString || "";
    return addClassesAndReturnUndo(missingClasses(classString));
  }
  function setClassesFromObject(el, classObject) {
    let split = (classString) => classString.split(" ").filter(Boolean);
    let forAdd = Object.entries(classObject).flatMap(([classString, bool]) => bool ? split(classString) : false).filter(Boolean);
    let forRemove = Object.entries(classObject).flatMap(([classString, bool]) => !bool ? split(classString) : false).filter(Boolean);
    let added = [];
    let removed = [];
    forRemove.forEach((i) => {
      if (el.classList.contains(i)) {
        el.classList.remove(i);
        removed.push(i);
      }
    });
    forAdd.forEach((i) => {
      if (!el.classList.contains(i)) {
        el.classList.add(i);
        added.push(i);
      }
    });
    return () => {
      removed.forEach((i) => el.classList.add(i));
      added.forEach((i) => el.classList.remove(i));
    };
  }
  function setStyles(el, value) {
    if (typeof value === "object" && value !== null) {
      return setStylesFromObject(el, value);
    }
    return setStylesFromString(el, value);
  }
  function setStylesFromObject(el, value) {
    let previousStyles = {};
    Object.entries(value).forEach(([key, value2]) => {
      previousStyles[key] = el.style[key];
      if (!key.startsWith("--")) {
        key = kebabCase(key);
      }
      el.style.setProperty(key, value2);
    });
    setTimeout(() => {
      if (el.style.length === 0) {
        el.removeAttribute("style");
      }
    });
    return () => {
      setStyles(el, previousStyles);
    };
  }
  function setStylesFromString(el, value) {
    let cache = el.getAttribute("style", value);
    el.setAttribute("style", value);
    return () => {
      el.setAttribute("style", cache || "");
    };
  }
  function kebabCase(subject) {
    return subject.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
  }
  function once(callback, fallback = () => {
  }) {
    let called = false;
    return function() {
      if (!called) {
        called = true;
        callback.apply(this, arguments);
      } else {
        fallback.apply(this, arguments);
      }
    };
  }
  directive("transition", (el, { value, modifiers, expression }, { evaluate: evaluate2 }) => {
    if (typeof expression === "function")
      expression = evaluate2(expression);
    if (expression === false)
      return;
    if (!expression || typeof expression === "boolean") {
      registerTransitionsFromHelper(el, modifiers, value);
    } else {
      registerTransitionsFromClassString(el, expression, value);
    }
  });
  function registerTransitionsFromClassString(el, classString, stage) {
    registerTransitionObject(el, setClasses, "");
    let directiveStorageMap = {
      "enter": (classes) => {
        el._x_transition.enter.during = classes;
      },
      "enter-start": (classes) => {
        el._x_transition.enter.start = classes;
      },
      "enter-end": (classes) => {
        el._x_transition.enter.end = classes;
      },
      "leave": (classes) => {
        el._x_transition.leave.during = classes;
      },
      "leave-start": (classes) => {
        el._x_transition.leave.start = classes;
      },
      "leave-end": (classes) => {
        el._x_transition.leave.end = classes;
      }
    };
    directiveStorageMap[stage](classString);
  }
  function registerTransitionsFromHelper(el, modifiers, stage) {
    registerTransitionObject(el, setStyles);
    let doesntSpecify = !modifiers.includes("in") && !modifiers.includes("out") && !stage;
    let transitioningIn = doesntSpecify || modifiers.includes("in") || ["enter"].includes(stage);
    let transitioningOut = doesntSpecify || modifiers.includes("out") || ["leave"].includes(stage);
    if (modifiers.includes("in") && !doesntSpecify) {
      modifiers = modifiers.filter((i, index) => index < modifiers.indexOf("out"));
    }
    if (modifiers.includes("out") && !doesntSpecify) {
      modifiers = modifiers.filter((i, index) => index > modifiers.indexOf("out"));
    }
    let wantsAll = !modifiers.includes("opacity") && !modifiers.includes("scale");
    let wantsOpacity = wantsAll || modifiers.includes("opacity");
    let wantsScale = wantsAll || modifiers.includes("scale");
    let opacityValue = wantsOpacity ? 0 : 1;
    let scaleValue = wantsScale ? modifierValue(modifiers, "scale", 95) / 100 : 1;
    let delay = modifierValue(modifiers, "delay", 0) / 1e3;
    let origin = modifierValue(modifiers, "origin", "center");
    let property = "opacity, transform";
    let durationIn = modifierValue(modifiers, "duration", 150) / 1e3;
    let durationOut = modifierValue(modifiers, "duration", 75) / 1e3;
    let easing = `cubic-bezier(0.4, 0.0, 0.2, 1)`;
    if (transitioningIn) {
      el._x_transition.enter.during = {
        transformOrigin: origin,
        transitionDelay: `${delay}s`,
        transitionProperty: property,
        transitionDuration: `${durationIn}s`,
        transitionTimingFunction: easing
      };
      el._x_transition.enter.start = {
        opacity: opacityValue,
        transform: `scale(${scaleValue})`
      };
      el._x_transition.enter.end = {
        opacity: 1,
        transform: `scale(1)`
      };
    }
    if (transitioningOut) {
      el._x_transition.leave.during = {
        transformOrigin: origin,
        transitionDelay: `${delay}s`,
        transitionProperty: property,
        transitionDuration: `${durationOut}s`,
        transitionTimingFunction: easing
      };
      el._x_transition.leave.start = {
        opacity: 1,
        transform: `scale(1)`
      };
      el._x_transition.leave.end = {
        opacity: opacityValue,
        transform: `scale(${scaleValue})`
      };
    }
  }
  function registerTransitionObject(el, setFunction, defaultValue = {}) {
    if (!el._x_transition)
      el._x_transition = {
        enter: { during: defaultValue, start: defaultValue, end: defaultValue },
        leave: { during: defaultValue, start: defaultValue, end: defaultValue },
        in(before = () => {
        }, after = () => {
        }) {
          transition(el, setFunction, {
            during: this.enter.during,
            start: this.enter.start,
            end: this.enter.end
          }, before, after);
        },
        out(before = () => {
        }, after = () => {
        }) {
          transition(el, setFunction, {
            during: this.leave.during,
            start: this.leave.start,
            end: this.leave.end
          }, before, after);
        }
      };
  }
  window.Element.prototype._x_toggleAndCascadeWithTransitions = function(el, value, show, hide) {
    const nextTick2 = document.visibilityState === "visible" ? requestAnimationFrame : setTimeout;
    let clickAwayCompatibleShow = () => nextTick2(show);
    if (value) {
      if (el._x_transition && (el._x_transition.enter || el._x_transition.leave)) {
        el._x_transition.enter && (Object.entries(el._x_transition.enter.during).length || Object.entries(el._x_transition.enter.start).length || Object.entries(el._x_transition.enter.end).length) ? el._x_transition.in(show) : clickAwayCompatibleShow();
      } else {
        el._x_transition ? el._x_transition.in(show) : clickAwayCompatibleShow();
      }
      return;
    }
    el._x_hidePromise = el._x_transition ? new Promise((resolve, reject) => {
      el._x_transition.out(() => {
      }, () => resolve(hide));
      el._x_transitioning && el._x_transitioning.beforeCancel(() => reject({ isFromCancelledTransition: true }));
    }) : Promise.resolve(hide);
    queueMicrotask(() => {
      let closest = closestHide(el);
      if (closest) {
        if (!closest._x_hideChildren)
          closest._x_hideChildren = [];
        closest._x_hideChildren.push(el);
      } else {
        nextTick2(() => {
          let hideAfterChildren = (el2) => {
            let carry = Promise.all([
              el2._x_hidePromise,
              ...(el2._x_hideChildren || []).map(hideAfterChildren)
            ]).then(([i]) => i?.());
            delete el2._x_hidePromise;
            delete el2._x_hideChildren;
            return carry;
          };
          hideAfterChildren(el).catch((e) => {
            if (!e.isFromCancelledTransition)
              throw e;
          });
        });
      }
    });
  };
  function closestHide(el) {
    let parent = el.parentNode;
    if (!parent)
      return;
    return parent._x_hidePromise ? parent : closestHide(parent);
  }
  function transition(el, setFunction, { during, start: start2, end } = {}, before = () => {
  }, after = () => {
  }) {
    if (el._x_transitioning)
      el._x_transitioning.cancel();
    if (Object.keys(during).length === 0 && Object.keys(start2).length === 0 && Object.keys(end).length === 0) {
      before();
      after();
      return;
    }
    let undoStart, undoDuring, undoEnd;
    performTransition(el, {
      start() {
        undoStart = setFunction(el, start2);
      },
      during() {
        undoDuring = setFunction(el, during);
      },
      before,
      end() {
        undoStart();
        undoEnd = setFunction(el, end);
      },
      after,
      cleanup() {
        undoDuring();
        undoEnd();
      }
    });
  }
  function performTransition(el, stages) {
    let interrupted, reachedBefore, reachedEnd;
    let finish = once(() => {
      mutateDom(() => {
        interrupted = true;
        if (!reachedBefore)
          stages.before();
        if (!reachedEnd) {
          stages.end();
          releaseNextTicks();
        }
        stages.after();
        if (el.isConnected)
          stages.cleanup();
        delete el._x_transitioning;
      });
    });
    el._x_transitioning = {
      beforeCancels: [],
      beforeCancel(callback) {
        this.beforeCancels.push(callback);
      },
      cancel: once(function() {
        while (this.beforeCancels.length) {
          this.beforeCancels.shift()();
        }
        ;
        finish();
      }),
      finish
    };
    mutateDom(() => {
      stages.start();
      stages.during();
    });
    holdNextTicks();
    requestAnimationFrame(() => {
      if (interrupted)
        return;
      let duration = Number(getComputedStyle(el).transitionDuration.replace(/,.*/, "").replace("s", "")) * 1e3;
      let delay = Number(getComputedStyle(el).transitionDelay.replace(/,.*/, "").replace("s", "")) * 1e3;
      if (duration === 0)
        duration = Number(getComputedStyle(el).animationDuration.replace("s", "")) * 1e3;
      mutateDom(() => {
        stages.before();
      });
      reachedBefore = true;
      requestAnimationFrame(() => {
        if (interrupted)
          return;
        mutateDom(() => {
          stages.end();
        });
        releaseNextTicks();
        setTimeout(el._x_transitioning.finish, duration + delay);
        reachedEnd = true;
      });
    });
  }
  function modifierValue(modifiers, key, fallback) {
    if (modifiers.indexOf(key) === -1)
      return fallback;
    const rawValue = modifiers[modifiers.indexOf(key) + 1];
    if (!rawValue)
      return fallback;
    if (key === "scale") {
      if (isNaN(rawValue))
        return fallback;
    }
    if (key === "duration" || key === "delay") {
      let match = rawValue.match(/([0-9]+)ms/);
      if (match)
        return match[1];
    }
    if (key === "origin") {
      if (["top", "right", "left", "center", "bottom"].includes(modifiers[modifiers.indexOf(key) + 2])) {
        return [rawValue, modifiers[modifiers.indexOf(key) + 2]].join(" ");
      }
    }
    return rawValue;
  }
  var isCloning = false;
  function skipDuringClone(callback, fallback = () => {
  }) {
    return (...args) => isCloning ? fallback(...args) : callback(...args);
  }
  function onlyDuringClone(callback) {
    return (...args) => isCloning && callback(...args);
  }
  var interceptors = [];
  function interceptClone(callback) {
    interceptors.push(callback);
  }
  function cloneNode(from, to) {
    interceptors.forEach((i) => i(from, to));
    isCloning = true;
    dontRegisterReactiveSideEffects(() => {
      initTree(to, (el, callback) => {
        callback(el, () => {
        });
      });
    });
    isCloning = false;
  }
  var isCloningLegacy = false;
  function clone(oldEl, newEl) {
    if (!newEl._x_dataStack)
      newEl._x_dataStack = oldEl._x_dataStack;
    isCloning = true;
    isCloningLegacy = true;
    dontRegisterReactiveSideEffects(() => {
      cloneTree(newEl);
    });
    isCloning = false;
    isCloningLegacy = false;
  }
  function cloneTree(el) {
    let hasRunThroughFirstEl = false;
    let shallowWalker = (el2, callback) => {
      walk(el2, (el3, skip) => {
        if (hasRunThroughFirstEl && isRoot(el3))
          return skip();
        hasRunThroughFirstEl = true;
        callback(el3, skip);
      });
    };
    initTree(el, shallowWalker);
  }
  function dontRegisterReactiveSideEffects(callback) {
    let cache = effect;
    overrideEffect((callback2, el) => {
      let storedEffect = cache(callback2);
      release(storedEffect);
      return () => {
      };
    });
    callback();
    overrideEffect(cache);
  }
  function bind(el, name, value, modifiers = []) {
    if (!el._x_bindings)
      el._x_bindings = reactive({});
    el._x_bindings[name] = value;
    name = modifiers.includes("camel") ? camelCase(name) : name;
    switch (name) {
      case "value":
        bindInputValue(el, value);
        break;
      case "style":
        bindStyles(el, value);
        break;
      case "class":
        bindClasses(el, value);
        break;
      case "selected":
      case "checked":
        bindAttributeAndProperty(el, name, value);
        break;
      default:
        bindAttribute(el, name, value);
        break;
    }
  }
  function bindInputValue(el, value) {
    if (el.type === "radio") {
      if (el.attributes.value === void 0) {
        el.value = value;
      }
      if (window.fromModel) {
        if (typeof value === "boolean") {
          el.checked = safeParseBoolean(el.value) === value;
        } else {
          el.checked = checkedAttrLooseCompare(el.value, value);
        }
      }
    } else if (el.type === "checkbox") {
      if (Number.isInteger(value)) {
        el.value = value;
      } else if (!Array.isArray(value) && typeof value !== "boolean" && ![null, void 0].includes(value)) {
        el.value = String(value);
      } else {
        if (Array.isArray(value)) {
          el.checked = value.some((val) => checkedAttrLooseCompare(val, el.value));
        } else {
          el.checked = !!value;
        }
      }
    } else if (el.tagName === "SELECT") {
      updateSelect(el, value);
    } else {
      if (el.value === value)
        return;
      el.value = value === void 0 ? "" : value;
    }
  }
  function bindClasses(el, value) {
    if (el._x_undoAddedClasses)
      el._x_undoAddedClasses();
    el._x_undoAddedClasses = setClasses(el, value);
  }
  function bindStyles(el, value) {
    if (el._x_undoAddedStyles)
      el._x_undoAddedStyles();
    el._x_undoAddedStyles = setStyles(el, value);
  }
  function bindAttributeAndProperty(el, name, value) {
    bindAttribute(el, name, value);
    setPropertyIfChanged(el, name, value);
  }
  function bindAttribute(el, name, value) {
    if ([null, void 0, false].includes(value) && attributeShouldntBePreservedIfFalsy(name)) {
      el.removeAttribute(name);
    } else {
      if (isBooleanAttr(name))
        value = name;
      setIfChanged(el, name, value);
    }
  }
  function setIfChanged(el, attrName, value) {
    if (el.getAttribute(attrName) != value) {
      el.setAttribute(attrName, value);
    }
  }
  function setPropertyIfChanged(el, propName, value) {
    if (el[propName] !== value) {
      el[propName] = value;
    }
  }
  function updateSelect(el, value) {
    const arrayWrappedValue = [].concat(value).map((value2) => {
      return value2 + "";
    });
    Array.from(el.options).forEach((option) => {
      option.selected = arrayWrappedValue.includes(option.value);
    });
  }
  function camelCase(subject) {
    return subject.toLowerCase().replace(/-(\w)/g, (match, char) => char.toUpperCase());
  }
  function checkedAttrLooseCompare(valueA, valueB) {
    return valueA == valueB;
  }
  function safeParseBoolean(rawValue) {
    if ([1, "1", "true", "on", "yes", true].includes(rawValue)) {
      return true;
    }
    if ([0, "0", "false", "off", "no", false].includes(rawValue)) {
      return false;
    }
    return rawValue ? Boolean(rawValue) : null;
  }
  function isBooleanAttr(attrName) {
    const booleanAttributes = [
      "disabled",
      "checked",
      "required",
      "readonly",
      "open",
      "selected",
      "autofocus",
      "itemscope",
      "multiple",
      "novalidate",
      "allowfullscreen",
      "allowpaymentrequest",
      "formnovalidate",
      "autoplay",
      "controls",
      "loop",
      "muted",
      "playsinline",
      "default",
      "ismap",
      "reversed",
      "async",
      "defer",
      "nomodule"
    ];
    return booleanAttributes.includes(attrName);
  }
  function attributeShouldntBePreservedIfFalsy(name) {
    return !["aria-pressed", "aria-checked", "aria-expanded", "aria-selected"].includes(name);
  }
  function getBinding(el, name, fallback) {
    if (el._x_bindings && el._x_bindings[name] !== void 0)
      return el._x_bindings[name];
    return getAttributeBinding(el, name, fallback);
  }
  function extractProp(el, name, fallback, extract = true) {
    if (el._x_bindings && el._x_bindings[name] !== void 0)
      return el._x_bindings[name];
    if (el._x_inlineBindings && el._x_inlineBindings[name] !== void 0) {
      let binding = el._x_inlineBindings[name];
      binding.extract = extract;
      return dontAutoEvaluateFunctions(() => {
        return evaluate(el, binding.expression);
      });
    }
    return getAttributeBinding(el, name, fallback);
  }
  function getAttributeBinding(el, name, fallback) {
    let attr = el.getAttribute(name);
    if (attr === null)
      return typeof fallback === "function" ? fallback() : fallback;
    if (attr === "")
      return true;
    if (isBooleanAttr(name)) {
      return !![name, "true"].includes(attr);
    }
    return attr;
  }
  function debounce(func, wait) {
    var timeout;
    return function() {
      var context = this, args = arguments;
      var later = function() {
        timeout = null;
        func.apply(context, args);
      };
      clearTimeout(timeout);
      timeout = setTimeout(later, wait);
    };
  }
  function throttle(func, limit) {
    let inThrottle;
    return function() {
      let context = this, args = arguments;
      if (!inThrottle) {
        func.apply(context, args);
        inThrottle = true;
        setTimeout(() => inThrottle = false, limit);
      }
    };
  }
  function entangle({ get: outerGet, set: outerSet }, { get: innerGet, set: innerSet }) {
    let firstRun = true;
    let outerHash;
    let innerHash;
    let reference = effect(() => {
      let outer = outerGet();
      let inner = innerGet();
      if (firstRun) {
        innerSet(cloneIfObject(outer));
        firstRun = false;
      } else {
        let outerHashLatest = JSON.stringify(outer);
        let innerHashLatest = JSON.stringify(inner);
        if (outerHashLatest !== outerHash) {
          innerSet(cloneIfObject(outer));
        } else if (outerHashLatest !== innerHashLatest) {
          outerSet(cloneIfObject(inner));
        } else {
        }
      }
      outerHash = JSON.stringify(outerGet());
      innerHash = JSON.stringify(innerGet());
    });
    return () => {
      release(reference);
    };
  }
  function cloneIfObject(value) {
    return typeof value === "object" ? JSON.parse(JSON.stringify(value)) : value;
  }
  function plugin(callback) {
    let callbacks = Array.isArray(callback) ? callback : [callback];
    callbacks.forEach((i) => i(alpine_default));
  }
  var stores = {};
  var isReactive = false;
  function store(name, value) {
    if (!isReactive) {
      stores = reactive(stores);
      isReactive = true;
    }
    if (value === void 0) {
      return stores[name];
    }
    stores[name] = value;
    if (typeof value === "object" && value !== null && value.hasOwnProperty("init") && typeof value.init === "function") {
      stores[name].init();
    }
    initInterceptors(stores[name]);
  }
  function getStores() {
    return stores;
  }
  var binds = {};
  function bind2(name, bindings) {
    let getBindings = typeof bindings !== "function" ? () => bindings : bindings;
    if (name instanceof Element) {
      return applyBindingsObject(name, getBindings());
    } else {
      binds[name] = getBindings;
    }
    return () => {
    };
  }
  function injectBindingProviders(obj) {
    Object.entries(binds).forEach(([name, callback]) => {
      Object.defineProperty(obj, name, {
        get() {
          return (...args) => {
            return callback(...args);
          };
        }
      });
    });
    return obj;
  }
  function applyBindingsObject(el, obj, original) {
    let cleanupRunners = [];
    while (cleanupRunners.length)
      cleanupRunners.pop()();
    let attributes = Object.entries(obj).map(([name, value]) => ({ name, value }));
    let staticAttributes = attributesOnly(attributes);
    attributes = attributes.map((attribute) => {
      if (staticAttributes.find((attr) => attr.name === attribute.name)) {
        return {
          name: `x-bind:${attribute.name}`,
          value: `"${attribute.value}"`
        };
      }
      return attribute;
    });
    directives(el, attributes, original).map((handle) => {
      cleanupRunners.push(handle.runCleanups);
      handle();
    });
    return () => {
      while (cleanupRunners.length)
        cleanupRunners.pop()();
    };
  }
  var datas = {};
  function data(name, callback) {
    datas[name] = callback;
  }
  function injectDataProviders(obj, context) {
    Object.entries(datas).forEach(([name, callback]) => {
      Object.defineProperty(obj, name, {
        get() {
          return (...args) => {
            return callback.bind(context)(...args);
          };
        },
        enumerable: false
      });
    });
    return obj;
  }
  var Alpine = {
    get reactive() {
      return reactive;
    },
    get release() {
      return release;
    },
    get effect() {
      return effect;
    },
    get raw() {
      return raw;
    },
    version: "3.14.1",
    flushAndStopDeferringMutations,
    dontAutoEvaluateFunctions,
    disableEffectScheduling,
    startObservingMutations,
    stopObservingMutations,
    setReactivityEngine,
    onAttributeRemoved,
    onAttributesAdded,
    closestDataStack,
    skipDuringClone,
    onlyDuringClone,
    addRootSelector,
    addInitSelector,
    interceptClone,
    addScopeToNode,
    deferMutations,
    mapAttributes,
    evaluateLater,
    interceptInit,
    setEvaluator,
    mergeProxies,
    extractProp,
    findClosest,
    onElRemoved,
    closestRoot,
    destroyTree,
    interceptor,
    // INTERNAL: not public API and is subject to change without major release.
    transition,
    // INTERNAL
    setStyles,
    // INTERNAL
    mutateDom,
    directive,
    entangle,
    throttle,
    debounce,
    evaluate,
    initTree,
    nextTick,
    prefixed: prefix,
    prefix: setPrefix,
    plugin,
    magic,
    store,
    start,
    clone,
    // INTERNAL
    cloneNode,
    // INTERNAL
    bound: getBinding,
    $data: scope,
    watch,
    walk,
    data,
    bind: bind2
  };
  var alpine_default = Alpine;
  function makeMap(str, expectsLowerCase) {
    const map = /* @__PURE__ */ Object.create(null);
    const list = str.split(",");
    for (let i = 0; i < list.length; i++) {
      map[list[i]] = true;
    }
    return expectsLowerCase ? (val) => !!map[val.toLowerCase()] : (val) => !!map[val];
  }
  var specialBooleanAttrs = `itemscope,allowfullscreen,formnovalidate,ismap,nomodule,novalidate,readonly`;
  var isBooleanAttr2 = /* @__PURE__ */ makeMap(specialBooleanAttrs + `,async,autofocus,autoplay,controls,default,defer,disabled,hidden,loop,open,required,reversed,scoped,seamless,checked,muted,multiple,selected`);
  var EMPTY_OBJ = true ? Object.freeze({}) : {};
  var EMPTY_ARR = true ? Object.freeze([]) : [];
  var hasOwnProperty = Object.prototype.hasOwnProperty;
  var hasOwn = (val, key) => hasOwnProperty.call(val, key);
  var isArray = Array.isArray;
  var isMap = (val) => toTypeString(val) === "[object Map]";
  var isString = (val) => typeof val === "string";
  var isSymbol = (val) => typeof val === "symbol";
  var isObject = (val) => val !== null && typeof val === "object";
  var objectToString = Object.prototype.toString;
  var toTypeString = (value) => objectToString.call(value);
  var toRawType = (value) => {
    return toTypeString(value).slice(8, -1);
  };
  var isIntegerKey = (key) => isString(key) && key !== "NaN" && key[0] !== "-" && "" + parseInt(key, 10) === key;
  var cacheStringFunction = (fn) => {
    const cache = /* @__PURE__ */ Object.create(null);
    return (str) => {
      const hit = cache[str];
      return hit || (cache[str] = fn(str));
    };
  };
  var camelizeRE = /-(\w)/g;
  var camelize = cacheStringFunction((str) => {
    return str.replace(camelizeRE, (_, c) => c ? c.toUpperCase() : "");
  });
  var hyphenateRE = /\B([A-Z])/g;
  var hyphenate = cacheStringFunction((str) => str.replace(hyphenateRE, "-$1").toLowerCase());
  var capitalize = cacheStringFunction((str) => str.charAt(0).toUpperCase() + str.slice(1));
  var toHandlerKey = cacheStringFunction((str) => str ? `on${capitalize(str)}` : ``);
  var hasChanged = (value, oldValue) => value !== oldValue && (value === value || oldValue === oldValue);
  var targetMap = /* @__PURE__ */ new WeakMap();
  var effectStack = [];
  var activeEffect;
  var ITERATE_KEY = Symbol(true ? "iterate" : "");
  var MAP_KEY_ITERATE_KEY = Symbol(true ? "Map key iterate" : "");
  function isEffect(fn) {
    return fn && fn._isEffect === true;
  }
  function effect2(fn, options = EMPTY_OBJ) {
    if (isEffect(fn)) {
      fn = fn.raw;
    }
    const effect3 = createReactiveEffect(fn, options);
    if (!options.lazy) {
      effect3();
    }
    return effect3;
  }
  function stop(effect3) {
    if (effect3.active) {
      cleanup(effect3);
      if (effect3.options.onStop) {
        effect3.options.onStop();
      }
      effect3.active = false;
    }
  }
  var uid = 0;
  function createReactiveEffect(fn, options) {
    const effect3 = function reactiveEffect() {
      if (!effect3.active) {
        return fn();
      }
      if (!effectStack.includes(effect3)) {
        cleanup(effect3);
        try {
          enableTracking();
          effectStack.push(effect3);
          activeEffect = effect3;
          return fn();
        } finally {
          effectStack.pop();
          resetTracking();
          activeEffect = effectStack[effectStack.length - 1];
        }
      }
    };
    effect3.id = uid++;
    effect3.allowRecurse = !!options.allowRecurse;
    effect3._isEffect = true;
    effect3.active = true;
    effect3.raw = fn;
    effect3.deps = [];
    effect3.options = options;
    return effect3;
  }
  function cleanup(effect3) {
    const { deps } = effect3;
    if (deps.length) {
      for (let i = 0; i < deps.length; i++) {
        deps[i].delete(effect3);
      }
      deps.length = 0;
    }
  }
  var shouldTrack = true;
  var trackStack = [];
  function pauseTracking() {
    trackStack.push(shouldTrack);
    shouldTrack = false;
  }
  function enableTracking() {
    trackStack.push(shouldTrack);
    shouldTrack = true;
  }
  function resetTracking() {
    const last = trackStack.pop();
    shouldTrack = last === void 0 ? true : last;
  }
  function track(target, type, key) {
    if (!shouldTrack || activeEffect === void 0) {
      return;
    }
    let depsMap = targetMap.get(target);
    if (!depsMap) {
      targetMap.set(target, depsMap = /* @__PURE__ */ new Map());
    }
    let dep = depsMap.get(key);
    if (!dep) {
      depsMap.set(key, dep = /* @__PURE__ */ new Set());
    }
    if (!dep.has(activeEffect)) {
      dep.add(activeEffect);
      activeEffect.deps.push(dep);
      if (activeEffect.options.onTrack) {
        activeEffect.options.onTrack({
          effect: activeEffect,
          target,
          type,
          key
        });
      }
    }
  }
  function trigger(target, type, key, newValue, oldValue, oldTarget) {
    const depsMap = targetMap.get(target);
    if (!depsMap) {
      return;
    }
    const effects = /* @__PURE__ */ new Set();
    const add2 = (effectsToAdd) => {
      if (effectsToAdd) {
        effectsToAdd.forEach((effect3) => {
          if (effect3 !== activeEffect || effect3.allowRecurse) {
            effects.add(effect3);
          }
        });
      }
    };
    if (type === "clear") {
      depsMap.forEach(add2);
    } else if (key === "length" && isArray(target)) {
      depsMap.forEach((dep, key2) => {
        if (key2 === "length" || key2 >= newValue) {
          add2(dep);
        }
      });
    } else {
      if (key !== void 0) {
        add2(depsMap.get(key));
      }
      switch (type) {
        case "add":
          if (!isArray(target)) {
            add2(depsMap.get(ITERATE_KEY));
            if (isMap(target)) {
              add2(depsMap.get(MAP_KEY_ITERATE_KEY));
            }
          } else if (isIntegerKey(key)) {
            add2(depsMap.get("length"));
          }
          break;
        case "delete":
          if (!isArray(target)) {
            add2(depsMap.get(ITERATE_KEY));
            if (isMap(target)) {
              add2(depsMap.get(MAP_KEY_ITERATE_KEY));
            }
          }
          break;
        case "set":
          if (isMap(target)) {
            add2(depsMap.get(ITERATE_KEY));
          }
          break;
      }
    }
    const run = (effect3) => {
      if (effect3.options.onTrigger) {
        effect3.options.onTrigger({
          effect: effect3,
          target,
          key,
          type,
          newValue,
          oldValue,
          oldTarget
        });
      }
      if (effect3.options.scheduler) {
        effect3.options.scheduler(effect3);
      } else {
        effect3();
      }
    };
    effects.forEach(run);
  }
  var isNonTrackableKeys = /* @__PURE__ */ makeMap(`__proto__,__v_isRef,__isVue`);
  var builtInSymbols = new Set(Object.getOwnPropertyNames(Symbol).map((key) => Symbol[key]).filter(isSymbol));
  var get2 = /* @__PURE__ */ createGetter();
  var readonlyGet = /* @__PURE__ */ createGetter(true);
  var arrayInstrumentations = /* @__PURE__ */ createArrayInstrumentations();
  function createArrayInstrumentations() {
    const instrumentations = {};
    ["includes", "indexOf", "lastIndexOf"].forEach((key) => {
      instrumentations[key] = function(...args) {
        const arr = toRaw(this);
        for (let i = 0, l = this.length; i < l; i++) {
          track(arr, "get", i + "");
        }
        const res = arr[key](...args);
        if (res === -1 || res === false) {
          return arr[key](...args.map(toRaw));
        } else {
          return res;
        }
      };
    });
    ["push", "pop", "shift", "unshift", "splice"].forEach((key) => {
      instrumentations[key] = function(...args) {
        pauseTracking();
        const res = toRaw(this)[key].apply(this, args);
        resetTracking();
        return res;
      };
    });
    return instrumentations;
  }
  function createGetter(isReadonly = false, shallow = false) {
    return function get3(target, key, receiver) {
      if (key === "__v_isReactive") {
        return !isReadonly;
      } else if (key === "__v_isReadonly") {
        return isReadonly;
      } else if (key === "__v_raw" && receiver === (isReadonly ? shallow ? shallowReadonlyMap : readonlyMap : shallow ? shallowReactiveMap : reactiveMap).get(target)) {
        return target;
      }
      const targetIsArray = isArray(target);
      if (!isReadonly && targetIsArray && hasOwn(arrayInstrumentations, key)) {
        return Reflect.get(arrayInstrumentations, key, receiver);
      }
      const res = Reflect.get(target, key, receiver);
      if (isSymbol(key) ? builtInSymbols.has(key) : isNonTrackableKeys(key)) {
        return res;
      }
      if (!isReadonly) {
        track(target, "get", key);
      }
      if (shallow) {
        return res;
      }
      if (isRef(res)) {
        const shouldUnwrap = !targetIsArray || !isIntegerKey(key);
        return shouldUnwrap ? res.value : res;
      }
      if (isObject(res)) {
        return isReadonly ? readonly(res) : reactive2(res);
      }
      return res;
    };
  }
  var set2 = /* @__PURE__ */ createSetter();
  function createSetter(shallow = false) {
    return function set3(target, key, value, receiver) {
      let oldValue = target[key];
      if (!shallow) {
        value = toRaw(value);
        oldValue = toRaw(oldValue);
        if (!isArray(target) && isRef(oldValue) && !isRef(value)) {
          oldValue.value = value;
          return true;
        }
      }
      const hadKey = isArray(target) && isIntegerKey(key) ? Number(key) < target.length : hasOwn(target, key);
      const result = Reflect.set(target, key, value, receiver);
      if (target === toRaw(receiver)) {
        if (!hadKey) {
          trigger(target, "add", key, value);
        } else if (hasChanged(value, oldValue)) {
          trigger(target, "set", key, value, oldValue);
        }
      }
      return result;
    };
  }
  function deleteProperty(target, key) {
    const hadKey = hasOwn(target, key);
    const oldValue = target[key];
    const result = Reflect.deleteProperty(target, key);
    if (result && hadKey) {
      trigger(target, "delete", key, void 0, oldValue);
    }
    return result;
  }
  function has(target, key) {
    const result = Reflect.has(target, key);
    if (!isSymbol(key) || !builtInSymbols.has(key)) {
      track(target, "has", key);
    }
    return result;
  }
  function ownKeys(target) {
    track(target, "iterate", isArray(target) ? "length" : ITERATE_KEY);
    return Reflect.ownKeys(target);
  }
  var mutableHandlers = {
    get: get2,
    set: set2,
    deleteProperty,
    has,
    ownKeys
  };
  var readonlyHandlers = {
    get: readonlyGet,
    set(target, key) {
      if (true) {
        console.warn(`Set operation on key "${String(key)}" failed: target is readonly.`, target);
      }
      return true;
    },
    deleteProperty(target, key) {
      if (true) {
        console.warn(`Delete operation on key "${String(key)}" failed: target is readonly.`, target);
      }
      return true;
    }
  };
  var toReactive = (value) => isObject(value) ? reactive2(value) : value;
  var toReadonly = (value) => isObject(value) ? readonly(value) : value;
  var toShallow = (value) => value;
  var getProto = (v) => Reflect.getPrototypeOf(v);
  function get$1(target, key, isReadonly = false, isShallow = false) {
    target = target[
      "__v_raw"
      /* RAW */
    ];
    const rawTarget = toRaw(target);
    const rawKey = toRaw(key);
    if (key !== rawKey) {
      !isReadonly && track(rawTarget, "get", key);
    }
    !isReadonly && track(rawTarget, "get", rawKey);
    const { has: has2 } = getProto(rawTarget);
    const wrap = isShallow ? toShallow : isReadonly ? toReadonly : toReactive;
    if (has2.call(rawTarget, key)) {
      return wrap(target.get(key));
    } else if (has2.call(rawTarget, rawKey)) {
      return wrap(target.get(rawKey));
    } else if (target !== rawTarget) {
      target.get(key);
    }
  }
  function has$1(key, isReadonly = false) {
    const target = this[
      "__v_raw"
      /* RAW */
    ];
    const rawTarget = toRaw(target);
    const rawKey = toRaw(key);
    if (key !== rawKey) {
      !isReadonly && track(rawTarget, "has", key);
    }
    !isReadonly && track(rawTarget, "has", rawKey);
    return key === rawKey ? target.has(key) : target.has(key) || target.has(rawKey);
  }
  function size(target, isReadonly = false) {
    target = target[
      "__v_raw"
      /* RAW */
    ];
    !isReadonly && track(toRaw(target), "iterate", ITERATE_KEY);
    return Reflect.get(target, "size", target);
  }
  function add(value) {
    value = toRaw(value);
    const target = toRaw(this);
    const proto = getProto(target);
    const hadKey = proto.has.call(target, value);
    if (!hadKey) {
      target.add(value);
      trigger(target, "add", value, value);
    }
    return this;
  }
  function set$1(key, value) {
    value = toRaw(value);
    const target = toRaw(this);
    const { has: has2, get: get3 } = getProto(target);
    let hadKey = has2.call(target, key);
    if (!hadKey) {
      key = toRaw(key);
      hadKey = has2.call(target, key);
    } else if (true) {
      checkIdentityKeys(target, has2, key);
    }
    const oldValue = get3.call(target, key);
    target.set(key, value);
    if (!hadKey) {
      trigger(target, "add", key, value);
    } else if (hasChanged(value, oldValue)) {
      trigger(target, "set", key, value, oldValue);
    }
    return this;
  }
  function deleteEntry(key) {
    const target = toRaw(this);
    const { has: has2, get: get3 } = getProto(target);
    let hadKey = has2.call(target, key);
    if (!hadKey) {
      key = toRaw(key);
      hadKey = has2.call(target, key);
    } else if (true) {
      checkIdentityKeys(target, has2, key);
    }
    const oldValue = get3 ? get3.call(target, key) : void 0;
    const result = target.delete(key);
    if (hadKey) {
      trigger(target, "delete", key, void 0, oldValue);
    }
    return result;
  }
  function clear() {
    const target = toRaw(this);
    const hadItems = target.size !== 0;
    const oldTarget = true ? isMap(target) ? new Map(target) : new Set(target) : void 0;
    const result = target.clear();
    if (hadItems) {
      trigger(target, "clear", void 0, void 0, oldTarget);
    }
    return result;
  }
  function createForEach(isReadonly, isShallow) {
    return function forEach(callback, thisArg) {
      const observed = this;
      const target = observed[
        "__v_raw"
        /* RAW */
      ];
      const rawTarget = toRaw(target);
      const wrap = isShallow ? toShallow : isReadonly ? toReadonly : toReactive;
      !isReadonly && track(rawTarget, "iterate", ITERATE_KEY);
      return target.forEach((value, key) => {
        return callback.call(thisArg, wrap(value), wrap(key), observed);
      });
    };
  }
  function createIterableMethod(method, isReadonly, isShallow) {
    return function(...args) {
      const target = this[
        "__v_raw"
        /* RAW */
      ];
      const rawTarget = toRaw(target);
      const targetIsMap = isMap(rawTarget);
      const isPair = method === "entries" || method === Symbol.iterator && targetIsMap;
      const isKeyOnly = method === "keys" && targetIsMap;
      const innerIterator = target[method](...args);
      const wrap = isShallow ? toShallow : isReadonly ? toReadonly : toReactive;
      !isReadonly && track(rawTarget, "iterate", isKeyOnly ? MAP_KEY_ITERATE_KEY : ITERATE_KEY);
      return {
        // iterator protocol
        next() {
          const { value, done } = innerIterator.next();
          return done ? { value, done } : {
            value: isPair ? [wrap(value[0]), wrap(value[1])] : wrap(value),
            done
          };
        },
        // iterable protocol
        [Symbol.iterator]() {
          return this;
        }
      };
    };
  }
  function createReadonlyMethod(type) {
    return function(...args) {
      if (true) {
        const key = args[0] ? `on key "${args[0]}" ` : ``;
        console.warn(`${capitalize(type)} operation ${key}failed: target is readonly.`, toRaw(this));
      }
      return type === "delete" ? false : this;
    };
  }
  function createInstrumentations() {
    const mutableInstrumentations2 = {
      get(key) {
        return get$1(this, key);
      },
      get size() {
        return size(this);
      },
      has: has$1,
      add,
      set: set$1,
      delete: deleteEntry,
      clear,
      forEach: createForEach(false, false)
    };
    const shallowInstrumentations2 = {
      get(key) {
        return get$1(this, key, false, true);
      },
      get size() {
        return size(this);
      },
      has: has$1,
      add,
      set: set$1,
      delete: deleteEntry,
      clear,
      forEach: createForEach(false, true)
    };
    const readonlyInstrumentations2 = {
      get(key) {
        return get$1(this, key, true);
      },
      get size() {
        return size(this, true);
      },
      has(key) {
        return has$1.call(this, key, true);
      },
      add: createReadonlyMethod(
        "add"
        /* ADD */
      ),
      set: createReadonlyMethod(
        "set"
        /* SET */
      ),
      delete: createReadonlyMethod(
        "delete"
        /* DELETE */
      ),
      clear: createReadonlyMethod(
        "clear"
        /* CLEAR */
      ),
      forEach: createForEach(true, false)
    };
    const shallowReadonlyInstrumentations2 = {
      get(key) {
        return get$1(this, key, true, true);
      },
      get size() {
        return size(this, true);
      },
      has(key) {
        return has$1.call(this, key, true);
      },
      add: createReadonlyMethod(
        "add"
        /* ADD */
      ),
      set: createReadonlyMethod(
        "set"
        /* SET */
      ),
      delete: createReadonlyMethod(
        "delete"
        /* DELETE */
      ),
      clear: createReadonlyMethod(
        "clear"
        /* CLEAR */
      ),
      forEach: createForEach(true, true)
    };
    const iteratorMethods = ["keys", "values", "entries", Symbol.iterator];
    iteratorMethods.forEach((method) => {
      mutableInstrumentations2[method] = createIterableMethod(method, false, false);
      readonlyInstrumentations2[method] = createIterableMethod(method, true, false);
      shallowInstrumentations2[method] = createIterableMethod(method, false, true);
      shallowReadonlyInstrumentations2[method] = createIterableMethod(method, true, true);
    });
    return [
      mutableInstrumentations2,
      readonlyInstrumentations2,
      shallowInstrumentations2,
      shallowReadonlyInstrumentations2
    ];
  }
  var [mutableInstrumentations, readonlyInstrumentations, shallowInstrumentations, shallowReadonlyInstrumentations] = /* @__PURE__ */ createInstrumentations();
  function createInstrumentationGetter(isReadonly, shallow) {
    const instrumentations = shallow ? isReadonly ? shallowReadonlyInstrumentations : shallowInstrumentations : isReadonly ? readonlyInstrumentations : mutableInstrumentations;
    return (target, key, receiver) => {
      if (key === "__v_isReactive") {
        return !isReadonly;
      } else if (key === "__v_isReadonly") {
        return isReadonly;
      } else if (key === "__v_raw") {
        return target;
      }
      return Reflect.get(hasOwn(instrumentations, key) && key in target ? instrumentations : target, key, receiver);
    };
  }
  var mutableCollectionHandlers = {
    get: /* @__PURE__ */ createInstrumentationGetter(false, false)
  };
  var readonlyCollectionHandlers = {
    get: /* @__PURE__ */ createInstrumentationGetter(true, false)
  };
  function checkIdentityKeys(target, has2, key) {
    const rawKey = toRaw(key);
    if (rawKey !== key && has2.call(target, rawKey)) {
      const type = toRawType(target);
      console.warn(`Reactive ${type} contains both the raw and reactive versions of the same object${type === `Map` ? ` as keys` : ``}, which can lead to inconsistencies. Avoid differentiating between the raw and reactive versions of an object and only use the reactive version if possible.`);
    }
  }
  var reactiveMap = /* @__PURE__ */ new WeakMap();
  var shallowReactiveMap = /* @__PURE__ */ new WeakMap();
  var readonlyMap = /* @__PURE__ */ new WeakMap();
  var shallowReadonlyMap = /* @__PURE__ */ new WeakMap();
  function targetTypeMap(rawType) {
    switch (rawType) {
      case "Object":
      case "Array":
        return 1;
      case "Map":
      case "Set":
      case "WeakMap":
      case "WeakSet":
        return 2;
      default:
        return 0;
    }
  }
  function getTargetType(value) {
    return value[
      "__v_skip"
      /* SKIP */
    ] || !Object.isExtensible(value) ? 0 : targetTypeMap(toRawType(value));
  }
  function reactive2(target) {
    if (target && target[
      "__v_isReadonly"
      /* IS_READONLY */
    ]) {
      return target;
    }
    return createReactiveObject(target, false, mutableHandlers, mutableCollectionHandlers, reactiveMap);
  }
  function readonly(target) {
    return createReactiveObject(target, true, readonlyHandlers, readonlyCollectionHandlers, readonlyMap);
  }
  function createReactiveObject(target, isReadonly, baseHandlers, collectionHandlers, proxyMap) {
    if (!isObject(target)) {
      if (true) {
        console.warn(`value cannot be made reactive: ${String(target)}`);
      }
      return target;
    }
    if (target[
      "__v_raw"
      /* RAW */
    ] && !(isReadonly && target[
      "__v_isReactive"
      /* IS_REACTIVE */
    ])) {
      return target;
    }
    const existingProxy = proxyMap.get(target);
    if (existingProxy) {
      return existingProxy;
    }
    const targetType = getTargetType(target);
    if (targetType === 0) {
      return target;
    }
    const proxy = new Proxy(target, targetType === 2 ? collectionHandlers : baseHandlers);
    proxyMap.set(target, proxy);
    return proxy;
  }
  function toRaw(observed) {
    return observed && toRaw(observed[
      "__v_raw"
      /* RAW */
    ]) || observed;
  }
  function isRef(r) {
    return Boolean(r && r.__v_isRef === true);
  }
  magic("nextTick", () => nextTick);
  magic("dispatch", (el) => dispatch.bind(dispatch, el));
  magic("watch", (el, { evaluateLater: evaluateLater2, cleanup: cleanup2 }) => (key, callback) => {
    let evaluate2 = evaluateLater2(key);
    let getter = () => {
      let value;
      evaluate2((i) => value = i);
      return value;
    };
    let unwatch = watch(getter, callback);
    cleanup2(unwatch);
  });
  magic("store", getStores);
  magic("data", (el) => scope(el));
  magic("root", (el) => closestRoot(el));
  magic("refs", (el) => {
    if (el._x_refs_proxy)
      return el._x_refs_proxy;
    el._x_refs_proxy = mergeProxies(getArrayOfRefObject(el));
    return el._x_refs_proxy;
  });
  function getArrayOfRefObject(el) {
    let refObjects = [];
    findClosest(el, (i) => {
      if (i._x_refs)
        refObjects.push(i._x_refs);
    });
    return refObjects;
  }
  var globalIdMemo = {};
  function findAndIncrementId(name) {
    if (!globalIdMemo[name])
      globalIdMemo[name] = 0;
    return ++globalIdMemo[name];
  }
  function closestIdRoot(el, name) {
    return findClosest(el, (element) => {
      if (element._x_ids && element._x_ids[name])
        return true;
    });
  }
  function setIdRoot(el, name) {
    if (!el._x_ids)
      el._x_ids = {};
    if (!el._x_ids[name])
      el._x_ids[name] = findAndIncrementId(name);
  }
  magic("id", (el, { cleanup: cleanup2 }) => (name, key = null) => {
    let cacheKey = `${name}${key ? `-${key}` : ""}`;
    return cacheIdByNameOnElement(el, cacheKey, cleanup2, () => {
      let root = closestIdRoot(el, name);
      let id = root ? root._x_ids[name] : findAndIncrementId(name);
      return key ? `${name}-${id}-${key}` : `${name}-${id}`;
    });
  });
  interceptClone((from, to) => {
    if (from._x_id) {
      to._x_id = from._x_id;
    }
  });
  function cacheIdByNameOnElement(el, cacheKey, cleanup2, callback) {
    if (!el._x_id)
      el._x_id = {};
    if (el._x_id[cacheKey])
      return el._x_id[cacheKey];
    let output = callback();
    el._x_id[cacheKey] = output;
    cleanup2(() => {
      delete el._x_id[cacheKey];
    });
    return output;
  }
  magic("el", (el) => el);
  warnMissingPluginMagic("Focus", "focus", "focus");
  warnMissingPluginMagic("Persist", "persist", "persist");
  function warnMissingPluginMagic(name, magicName, slug) {
    magic(magicName, (el) => warn(`You can't use [$${magicName}] without first installing the "${name}" plugin here: https://alpinejs.dev/plugins/${slug}`, el));
  }
  directive("modelable", (el, { expression }, { effect: effect3, evaluateLater: evaluateLater2, cleanup: cleanup2 }) => {
    let func = evaluateLater2(expression);
    let innerGet = () => {
      let result;
      func((i) => result = i);
      return result;
    };
    let evaluateInnerSet = evaluateLater2(`${expression} = __placeholder`);
    let innerSet = (val) => evaluateInnerSet(() => {
    }, { scope: { "__placeholder": val } });
    let initialValue = innerGet();
    innerSet(initialValue);
    queueMicrotask(() => {
      if (!el._x_model)
        return;
      el._x_removeModelListeners["default"]();
      let outerGet = el._x_model.get;
      let outerSet = el._x_model.set;
      let releaseEntanglement = entangle(
        {
          get() {
            return outerGet();
          },
          set(value) {
            outerSet(value);
          }
        },
        {
          get() {
            return innerGet();
          },
          set(value) {
            innerSet(value);
          }
        }
      );
      cleanup2(releaseEntanglement);
    });
  });
  directive("teleport", (el, { modifiers, expression }, { cleanup: cleanup2 }) => {
    if (el.tagName.toLowerCase() !== "template")
      warn("x-teleport can only be used on a <template> tag", el);
    let target = getTarget(expression);
    let clone2 = el.content.cloneNode(true).firstElementChild;
    el._x_teleport = clone2;
    clone2._x_teleportBack = el;
    el.setAttribute("data-teleport-template", true);
    clone2.setAttribute("data-teleport-target", true);
    if (el._x_forwardEvents) {
      el._x_forwardEvents.forEach((eventName) => {
        clone2.addEventListener(eventName, (e) => {
          e.stopPropagation();
          el.dispatchEvent(new e.constructor(e.type, e));
        });
      });
    }
    addScopeToNode(clone2, {}, el);
    let placeInDom = (clone3, target2, modifiers2) => {
      if (modifiers2.includes("prepend")) {
        target2.parentNode.insertBefore(clone3, target2);
      } else if (modifiers2.includes("append")) {
        target2.parentNode.insertBefore(clone3, target2.nextSibling);
      } else {
        target2.appendChild(clone3);
      }
    };
    mutateDom(() => {
      placeInDom(clone2, target, modifiers);
      skipDuringClone(() => {
        initTree(clone2);
        clone2._x_ignore = true;
      })();
    });
    el._x_teleportPutBack = () => {
      let target2 = getTarget(expression);
      mutateDom(() => {
        placeInDom(el._x_teleport, target2, modifiers);
      });
    };
    cleanup2(() => clone2.remove());
  });
  var teleportContainerDuringClone = document.createElement("div");
  function getTarget(expression) {
    let target = skipDuringClone(() => {
      return document.querySelector(expression);
    }, () => {
      return teleportContainerDuringClone;
    })();
    if (!target)
      warn(`Cannot find x-teleport element for selector: "${expression}"`);
    return target;
  }
  var handler = () => {
  };
  handler.inline = (el, { modifiers }, { cleanup: cleanup2 }) => {
    modifiers.includes("self") ? el._x_ignoreSelf = true : el._x_ignore = true;
    cleanup2(() => {
      modifiers.includes("self") ? delete el._x_ignoreSelf : delete el._x_ignore;
    });
  };
  directive("ignore", handler);
  directive("effect", skipDuringClone((el, { expression }, { effect: effect3 }) => {
    effect3(evaluateLater(el, expression));
  }));
  function on(el, event, modifiers, callback) {
    let listenerTarget = el;
    let handler4 = (e) => callback(e);
    let options = {};
    let wrapHandler = (callback2, wrapper) => (e) => wrapper(callback2, e);
    if (modifiers.includes("dot"))
      event = dotSyntax(event);
    if (modifiers.includes("camel"))
      event = camelCase2(event);
    if (modifiers.includes("passive"))
      options.passive = true;
    if (modifiers.includes("capture"))
      options.capture = true;
    if (modifiers.includes("window"))
      listenerTarget = window;
    if (modifiers.includes("document"))
      listenerTarget = document;
    if (modifiers.includes("debounce")) {
      let nextModifier = modifiers[modifiers.indexOf("debounce") + 1] || "invalid-wait";
      let wait = isNumeric(nextModifier.split("ms")[0]) ? Number(nextModifier.split("ms")[0]) : 250;
      handler4 = debounce(handler4, wait);
    }
    if (modifiers.includes("throttle")) {
      let nextModifier = modifiers[modifiers.indexOf("throttle") + 1] || "invalid-wait";
      let wait = isNumeric(nextModifier.split("ms")[0]) ? Number(nextModifier.split("ms")[0]) : 250;
      handler4 = throttle(handler4, wait);
    }
    if (modifiers.includes("prevent"))
      handler4 = wrapHandler(handler4, (next, e) => {
        e.preventDefault();
        next(e);
      });
    if (modifiers.includes("stop"))
      handler4 = wrapHandler(handler4, (next, e) => {
        e.stopPropagation();
        next(e);
      });
    if (modifiers.includes("once")) {
      handler4 = wrapHandler(handler4, (next, e) => {
        next(e);
        listenerTarget.removeEventListener(event, handler4, options);
      });
    }
    if (modifiers.includes("away") || modifiers.includes("outside")) {
      listenerTarget = document;
      handler4 = wrapHandler(handler4, (next, e) => {
        if (el.contains(e.target))
          return;
        if (e.target.isConnected === false)
          return;
        if (el.offsetWidth < 1 && el.offsetHeight < 1)
          return;
        if (el._x_isShown === false)
          return;
        next(e);
      });
    }
    if (modifiers.includes("self"))
      handler4 = wrapHandler(handler4, (next, e) => {
        e.target === el && next(e);
      });
    if (isKeyEvent(event) || isClickEvent(event)) {
      handler4 = wrapHandler(handler4, (next, e) => {
        if (isListeningForASpecificKeyThatHasntBeenPressed(e, modifiers)) {
          return;
        }
        next(e);
      });
    }
    listenerTarget.addEventListener(event, handler4, options);
    return () => {
      listenerTarget.removeEventListener(event, handler4, options);
    };
  }
  function dotSyntax(subject) {
    return subject.replace(/-/g, ".");
  }
  function camelCase2(subject) {
    return subject.toLowerCase().replace(/-(\w)/g, (match, char) => char.toUpperCase());
  }
  function isNumeric(subject) {
    return !Array.isArray(subject) && !isNaN(subject);
  }
  function kebabCase2(subject) {
    if ([" ", "_"].includes(
      subject
    ))
      return subject;
    return subject.replace(/([a-z])([A-Z])/g, "$1-$2").replace(/[_\s]/, "-").toLowerCase();
  }
  function isKeyEvent(event) {
    return ["keydown", "keyup"].includes(event);
  }
  function isClickEvent(event) {
    return ["contextmenu", "click", "mouse"].some((i) => event.includes(i));
  }
  function isListeningForASpecificKeyThatHasntBeenPressed(e, modifiers) {
    let keyModifiers = modifiers.filter((i) => {
      return !["window", "document", "prevent", "stop", "once", "capture", "self", "away", "outside", "passive"].includes(i);
    });
    if (keyModifiers.includes("debounce")) {
      let debounceIndex = keyModifiers.indexOf("debounce");
      keyModifiers.splice(debounceIndex, isNumeric((keyModifiers[debounceIndex + 1] || "invalid-wait").split("ms")[0]) ? 2 : 1);
    }
    if (keyModifiers.includes("throttle")) {
      let debounceIndex = keyModifiers.indexOf("throttle");
      keyModifiers.splice(debounceIndex, isNumeric((keyModifiers[debounceIndex + 1] || "invalid-wait").split("ms")[0]) ? 2 : 1);
    }
    if (keyModifiers.length === 0)
      return false;
    if (keyModifiers.length === 1 && keyToModifiers(e.key).includes(keyModifiers[0]))
      return false;
    const systemKeyModifiers = ["ctrl", "shift", "alt", "meta", "cmd", "super"];
    const selectedSystemKeyModifiers = systemKeyModifiers.filter((modifier) => keyModifiers.includes(modifier));
    keyModifiers = keyModifiers.filter((i) => !selectedSystemKeyModifiers.includes(i));
    if (selectedSystemKeyModifiers.length > 0) {
      const activelyPressedKeyModifiers = selectedSystemKeyModifiers.filter((modifier) => {
        if (modifier === "cmd" || modifier === "super")
          modifier = "meta";
        return e[`${modifier}Key`];
      });
      if (activelyPressedKeyModifiers.length === selectedSystemKeyModifiers.length) {
        if (isClickEvent(e.type))
          return false;
        if (keyToModifiers(e.key).includes(keyModifiers[0]))
          return false;
      }
    }
    return true;
  }
  function keyToModifiers(key) {
    if (!key)
      return [];
    key = kebabCase2(key);
    let modifierToKeyMap = {
      "ctrl": "control",
      "slash": "/",
      "space": " ",
      "spacebar": " ",
      "cmd": "meta",
      "esc": "escape",
      "up": "arrow-up",
      "down": "arrow-down",
      "left": "arrow-left",
      "right": "arrow-right",
      "period": ".",
      "comma": ",",
      "equal": "=",
      "minus": "-",
      "underscore": "_"
    };
    modifierToKeyMap[key] = key;
    return Object.keys(modifierToKeyMap).map((modifier) => {
      if (modifierToKeyMap[modifier] === key)
        return modifier;
    }).filter((modifier) => modifier);
  }
  directive("model", (el, { modifiers, expression }, { effect: effect3, cleanup: cleanup2 }) => {
    let scopeTarget = el;
    if (modifiers.includes("parent")) {
      scopeTarget = el.parentNode;
    }
    let evaluateGet = evaluateLater(scopeTarget, expression);
    let evaluateSet;
    if (typeof expression === "string") {
      evaluateSet = evaluateLater(scopeTarget, `${expression} = __placeholder`);
    } else if (typeof expression === "function" && typeof expression() === "string") {
      evaluateSet = evaluateLater(scopeTarget, `${expression()} = __placeholder`);
    } else {
      evaluateSet = () => {
      };
    }
    let getValue = () => {
      let result;
      evaluateGet((value) => result = value);
      return isGetterSetter(result) ? result.get() : result;
    };
    let setValue = (value) => {
      let result;
      evaluateGet((value2) => result = value2);
      if (isGetterSetter(result)) {
        result.set(value);
      } else {
        evaluateSet(() => {
        }, {
          scope: { "__placeholder": value }
        });
      }
    };
    if (typeof expression === "string" && el.type === "radio") {
      mutateDom(() => {
        if (!el.hasAttribute("name"))
          el.setAttribute("name", expression);
      });
    }
    var event = el.tagName.toLowerCase() === "select" || ["checkbox", "radio"].includes(el.type) || modifiers.includes("lazy") ? "change" : "input";
    let removeListener2 = isCloning ? () => {
    } : on(el, event, modifiers, (e) => {
      setValue(getInputValue(el, modifiers, e, getValue()));
    });
    if (modifiers.includes("fill")) {
      if ([void 0, null, ""].includes(getValue()) || el.type === "checkbox" && Array.isArray(getValue()) || el.tagName.toLowerCase() === "select" && el.multiple) {
        setValue(
          getInputValue(el, modifiers, { target: el }, getValue())
        );
      }
    }
    if (!el._x_removeModelListeners)
      el._x_removeModelListeners = {};
    el._x_removeModelListeners["default"] = removeListener2;
    cleanup2(() => el._x_removeModelListeners["default"]());
    if (el.form) {
      let removeResetListener = on(el.form, "reset", [], (e) => {
        nextTick(() => el._x_model && el._x_model.set(getInputValue(el, modifiers, { target: el }, getValue())));
      });
      cleanup2(() => removeResetListener());
    }
    el._x_model = {
      get() {
        return getValue();
      },
      set(value) {
        setValue(value);
      }
    };
    el._x_forceModelUpdate = (value) => {
      if (value === void 0 && typeof expression === "string" && expression.match(/\./))
        value = "";
      window.fromModel = true;
      mutateDom(() => bind(el, "value", value));
      delete window.fromModel;
    };
    effect3(() => {
      let value = getValue();
      if (modifiers.includes("unintrusive") && document.activeElement.isSameNode(el))
        return;
      el._x_forceModelUpdate(value);
    });
  });
  function getInputValue(el, modifiers, event, currentValue) {
    return mutateDom(() => {
      if (event instanceof CustomEvent && event.detail !== void 0)
        return event.detail !== null && event.detail !== void 0 ? event.detail : event.target.value;
      else if (el.type === "checkbox") {
        if (Array.isArray(currentValue)) {
          let newValue = null;
          if (modifiers.includes("number")) {
            newValue = safeParseNumber(event.target.value);
          } else if (modifiers.includes("boolean")) {
            newValue = safeParseBoolean(event.target.value);
          } else {
            newValue = event.target.value;
          }
          return event.target.checked ? currentValue.includes(newValue) ? currentValue : currentValue.concat([newValue]) : currentValue.filter((el2) => !checkedAttrLooseCompare2(el2, newValue));
        } else {
          return event.target.checked;
        }
      } else if (el.tagName.toLowerCase() === "select" && el.multiple) {
        if (modifiers.includes("number")) {
          return Array.from(event.target.selectedOptions).map((option) => {
            let rawValue = option.value || option.text;
            return safeParseNumber(rawValue);
          });
        } else if (modifiers.includes("boolean")) {
          return Array.from(event.target.selectedOptions).map((option) => {
            let rawValue = option.value || option.text;
            return safeParseBoolean(rawValue);
          });
        }
        return Array.from(event.target.selectedOptions).map((option) => {
          return option.value || option.text;
        });
      } else {
        let newValue;
        if (el.type === "radio") {
          if (event.target.checked) {
            newValue = event.target.value;
          } else {
            newValue = currentValue;
          }
        } else {
          newValue = event.target.value;
        }
        if (modifiers.includes("number")) {
          return safeParseNumber(newValue);
        } else if (modifiers.includes("boolean")) {
          return safeParseBoolean(newValue);
        } else if (modifiers.includes("trim")) {
          return newValue.trim();
        } else {
          return newValue;
        }
      }
    });
  }
  function safeParseNumber(rawValue) {
    let number = rawValue ? parseFloat(rawValue) : null;
    return isNumeric2(number) ? number : rawValue;
  }
  function checkedAttrLooseCompare2(valueA, valueB) {
    return valueA == valueB;
  }
  function isNumeric2(subject) {
    return !Array.isArray(subject) && !isNaN(subject);
  }
  function isGetterSetter(value) {
    return value !== null && typeof value === "object" && typeof value.get === "function" && typeof value.set === "function";
  }
  directive("cloak", (el) => queueMicrotask(() => mutateDom(() => el.removeAttribute(prefix("cloak")))));
  addInitSelector(() => `[${prefix("init")}]`);
  directive("init", skipDuringClone((el, { expression }, { evaluate: evaluate2 }) => {
    if (typeof expression === "string") {
      return !!expression.trim() && evaluate2(expression, {}, false);
    }
    return evaluate2(expression, {}, false);
  }));
  directive("text", (el, { expression }, { effect: effect3, evaluateLater: evaluateLater2 }) => {
    let evaluate2 = evaluateLater2(expression);
    effect3(() => {
      evaluate2((value) => {
        mutateDom(() => {
          el.textContent = value;
        });
      });
    });
  });
  directive("html", (el, { expression }, { effect: effect3, evaluateLater: evaluateLater2 }) => {
    let evaluate2 = evaluateLater2(expression);
    effect3(() => {
      evaluate2((value) => {
        mutateDom(() => {
          el.innerHTML = value;
          el._x_ignoreSelf = true;
          initTree(el);
          delete el._x_ignoreSelf;
        });
      });
    });
  });
  mapAttributes(startingWith(":", into(prefix("bind:"))));
  var handler2 = (el, { value, modifiers, expression, original }, { effect: effect3, cleanup: cleanup2 }) => {
    if (!value) {
      let bindingProviders = {};
      injectBindingProviders(bindingProviders);
      let getBindings = evaluateLater(el, expression);
      getBindings((bindings) => {
        applyBindingsObject(el, bindings, original);
      }, { scope: bindingProviders });
      return;
    }
    if (value === "key")
      return storeKeyForXFor(el, expression);
    if (el._x_inlineBindings && el._x_inlineBindings[value] && el._x_inlineBindings[value].extract) {
      return;
    }
    let evaluate2 = evaluateLater(el, expression);
    effect3(() => evaluate2((result) => {
      if (result === void 0 && typeof expression === "string" && expression.match(/\./)) {
        result = "";
      }
      mutateDom(() => bind(el, value, result, modifiers));
    }));
    cleanup2(() => {
      el._x_undoAddedClasses && el._x_undoAddedClasses();
      el._x_undoAddedStyles && el._x_undoAddedStyles();
    });
  };
  handler2.inline = (el, { value, modifiers, expression }) => {
    if (!value)
      return;
    if (!el._x_inlineBindings)
      el._x_inlineBindings = {};
    el._x_inlineBindings[value] = { expression, extract: false };
  };
  directive("bind", handler2);
  function storeKeyForXFor(el, expression) {
    el._x_keyExpression = expression;
  }
  addRootSelector(() => `[${prefix("data")}]`);
  directive("data", (el, { expression }, { cleanup: cleanup2 }) => {
    if (shouldSkipRegisteringDataDuringClone(el))
      return;
    expression = expression === "" ? "{}" : expression;
    let magicContext = {};
    injectMagics(magicContext, el);
    let dataProviderContext = {};
    injectDataProviders(dataProviderContext, magicContext);
    let data2 = evaluate(el, expression, { scope: dataProviderContext });
    if (data2 === void 0 || data2 === true)
      data2 = {};
    injectMagics(data2, el);
    let reactiveData = reactive(data2);
    initInterceptors(reactiveData);
    let undo = addScopeToNode(el, reactiveData);
    reactiveData["init"] && evaluate(el, reactiveData["init"]);
    cleanup2(() => {
      reactiveData["destroy"] && evaluate(el, reactiveData["destroy"]);
      undo();
    });
  });
  interceptClone((from, to) => {
    if (from._x_dataStack) {
      to._x_dataStack = from._x_dataStack;
      to.setAttribute("data-has-alpine-state", true);
    }
  });
  function shouldSkipRegisteringDataDuringClone(el) {
    if (!isCloning)
      return false;
    if (isCloningLegacy)
      return true;
    return el.hasAttribute("data-has-alpine-state");
  }
  directive("show", (el, { modifiers, expression }, { effect: effect3 }) => {
    let evaluate2 = evaluateLater(el, expression);
    if (!el._x_doHide)
      el._x_doHide = () => {
        mutateDom(() => {
          el.style.setProperty("display", "none", modifiers.includes("important") ? "important" : void 0);
        });
      };
    if (!el._x_doShow)
      el._x_doShow = () => {
        mutateDom(() => {
          if (el.style.length === 1 && el.style.display === "none") {
            el.removeAttribute("style");
          } else {
            el.style.removeProperty("display");
          }
        });
      };
    let hide = () => {
      el._x_doHide();
      el._x_isShown = false;
    };
    let show = () => {
      el._x_doShow();
      el._x_isShown = true;
    };
    let clickAwayCompatibleShow = () => setTimeout(show);
    let toggle = once(
      (value) => value ? show() : hide(),
      (value) => {
        if (typeof el._x_toggleAndCascadeWithTransitions === "function") {
          el._x_toggleAndCascadeWithTransitions(el, value, show, hide);
        } else {
          value ? clickAwayCompatibleShow() : hide();
        }
      }
    );
    let oldValue;
    let firstTime = true;
    effect3(() => evaluate2((value) => {
      if (!firstTime && value === oldValue)
        return;
      if (modifiers.includes("immediate"))
        value ? clickAwayCompatibleShow() : hide();
      toggle(value);
      oldValue = value;
      firstTime = false;
    }));
  });
  directive("for", (el, { expression }, { effect: effect3, cleanup: cleanup2 }) => {
    let iteratorNames = parseForExpression(expression);
    let evaluateItems = evaluateLater(el, iteratorNames.items);
    let evaluateKey = evaluateLater(
      el,
      // the x-bind:key expression is stored for our use instead of evaluated.
      el._x_keyExpression || "index"
    );
    el._x_prevKeys = [];
    el._x_lookup = {};
    effect3(() => loop(el, iteratorNames, evaluateItems, evaluateKey));
    cleanup2(() => {
      Object.values(el._x_lookup).forEach((el2) => el2.remove());
      delete el._x_prevKeys;
      delete el._x_lookup;
    });
  });
  function loop(el, iteratorNames, evaluateItems, evaluateKey) {
    let isObject2 = (i) => typeof i === "object" && !Array.isArray(i);
    let templateEl = el;
    evaluateItems((items) => {
      if (isNumeric3(items) && items >= 0) {
        items = Array.from(Array(items).keys(), (i) => i + 1);
      }
      if (items === void 0)
        items = [];
      let lookup = el._x_lookup;
      let prevKeys = el._x_prevKeys;
      let scopes = [];
      let keys = [];
      if (isObject2(items)) {
        items = Object.entries(items).map(([key, value]) => {
          let scope2 = getIterationScopeVariables(iteratorNames, value, key, items);
          evaluateKey((value2) => {
            if (keys.includes(value2))
              warn("Duplicate key on x-for", el);
            keys.push(value2);
          }, { scope: { index: key, ...scope2 } });
          scopes.push(scope2);
        });
      } else {
        for (let i = 0; i < items.length; i++) {
          let scope2 = getIterationScopeVariables(iteratorNames, items[i], i, items);
          evaluateKey((value) => {
            if (keys.includes(value))
              warn("Duplicate key on x-for", el);
            keys.push(value);
          }, { scope: { index: i, ...scope2 } });
          scopes.push(scope2);
        }
      }
      let adds = [];
      let moves = [];
      let removes = [];
      let sames = [];
      for (let i = 0; i < prevKeys.length; i++) {
        let key = prevKeys[i];
        if (keys.indexOf(key) === -1)
          removes.push(key);
      }
      prevKeys = prevKeys.filter((key) => !removes.includes(key));
      let lastKey = "template";
      for (let i = 0; i < keys.length; i++) {
        let key = keys[i];
        let prevIndex = prevKeys.indexOf(key);
        if (prevIndex === -1) {
          prevKeys.splice(i, 0, key);
          adds.push([lastKey, i]);
        } else if (prevIndex !== i) {
          let keyInSpot = prevKeys.splice(i, 1)[0];
          let keyForSpot = prevKeys.splice(prevIndex - 1, 1)[0];
          prevKeys.splice(i, 0, keyForSpot);
          prevKeys.splice(prevIndex, 0, keyInSpot);
          moves.push([keyInSpot, keyForSpot]);
        } else {
          sames.push(key);
        }
        lastKey = key;
      }
      for (let i = 0; i < removes.length; i++) {
        let key = removes[i];
        if (!!lookup[key]._x_effects) {
          lookup[key]._x_effects.forEach(dequeueJob);
        }
        lookup[key].remove();
        lookup[key] = null;
        delete lookup[key];
      }
      for (let i = 0; i < moves.length; i++) {
        let [keyInSpot, keyForSpot] = moves[i];
        let elInSpot = lookup[keyInSpot];
        let elForSpot = lookup[keyForSpot];
        let marker = document.createElement("div");
        mutateDom(() => {
          if (!elForSpot)
            warn(`x-for ":key" is undefined or invalid`, templateEl, keyForSpot, lookup);
          elForSpot.after(marker);
          elInSpot.after(elForSpot);
          elForSpot._x_currentIfEl && elForSpot.after(elForSpot._x_currentIfEl);
          marker.before(elInSpot);
          elInSpot._x_currentIfEl && elInSpot.after(elInSpot._x_currentIfEl);
          marker.remove();
        });
        elForSpot._x_refreshXForScope(scopes[keys.indexOf(keyForSpot)]);
      }
      for (let i = 0; i < adds.length; i++) {
        let [lastKey2, index] = adds[i];
        let lastEl = lastKey2 === "template" ? templateEl : lookup[lastKey2];
        if (lastEl._x_currentIfEl)
          lastEl = lastEl._x_currentIfEl;
        let scope2 = scopes[index];
        let key = keys[index];
        let clone2 = document.importNode(templateEl.content, true).firstElementChild;
        let reactiveScope = reactive(scope2);
        addScopeToNode(clone2, reactiveScope, templateEl);
        clone2._x_refreshXForScope = (newScope) => {
          Object.entries(newScope).forEach(([key2, value]) => {
            reactiveScope[key2] = value;
          });
        };
        mutateDom(() => {
          lastEl.after(clone2);
          skipDuringClone(() => initTree(clone2))();
        });
        if (typeof key === "object") {
          warn("x-for key cannot be an object, it must be a string or an integer", templateEl);
        }
        lookup[key] = clone2;
      }
      for (let i = 0; i < sames.length; i++) {
        lookup[sames[i]]._x_refreshXForScope(scopes[keys.indexOf(sames[i])]);
      }
      templateEl._x_prevKeys = keys;
    });
  }
  function parseForExpression(expression) {
    let forIteratorRE = /,([^,\}\]]*)(?:,([^,\}\]]*))?$/;
    let stripParensRE = /^\s*\(|\)\s*$/g;
    let forAliasRE = /([\s\S]*?)\s+(?:in|of)\s+([\s\S]*)/;
    let inMatch = expression.match(forAliasRE);
    if (!inMatch)
      return;
    let res = {};
    res.items = inMatch[2].trim();
    let item = inMatch[1].replace(stripParensRE, "").trim();
    let iteratorMatch = item.match(forIteratorRE);
    if (iteratorMatch) {
      res.item = item.replace(forIteratorRE, "").trim();
      res.index = iteratorMatch[1].trim();
      if (iteratorMatch[2]) {
        res.collection = iteratorMatch[2].trim();
      }
    } else {
      res.item = item;
    }
    return res;
  }
  function getIterationScopeVariables(iteratorNames, item, index, items) {
    let scopeVariables = {};
    if (/^\[.*\]$/.test(iteratorNames.item) && Array.isArray(item)) {
      let names = iteratorNames.item.replace("[", "").replace("]", "").split(",").map((i) => i.trim());
      names.forEach((name, i) => {
        scopeVariables[name] = item[i];
      });
    } else if (/^\{.*\}$/.test(iteratorNames.item) && !Array.isArray(item) && typeof item === "object") {
      let names = iteratorNames.item.replace("{", "").replace("}", "").split(",").map((i) => i.trim());
      names.forEach((name) => {
        scopeVariables[name] = item[name];
      });
    } else {
      scopeVariables[iteratorNames.item] = item;
    }
    if (iteratorNames.index)
      scopeVariables[iteratorNames.index] = index;
    if (iteratorNames.collection)
      scopeVariables[iteratorNames.collection] = items;
    return scopeVariables;
  }
  function isNumeric3(subject) {
    return !Array.isArray(subject) && !isNaN(subject);
  }
  function handler3() {
  }
  handler3.inline = (el, { expression }, { cleanup: cleanup2 }) => {
    let root = closestRoot(el);
    if (!root._x_refs)
      root._x_refs = {};
    root._x_refs[expression] = el;
    cleanup2(() => delete root._x_refs[expression]);
  };
  directive("ref", handler3);
  directive("if", (el, { expression }, { effect: effect3, cleanup: cleanup2 }) => {
    if (el.tagName.toLowerCase() !== "template")
      warn("x-if can only be used on a <template> tag", el);
    let evaluate2 = evaluateLater(el, expression);
    let show = () => {
      if (el._x_currentIfEl)
        return el._x_currentIfEl;
      let clone2 = el.content.cloneNode(true).firstElementChild;
      addScopeToNode(clone2, {}, el);
      mutateDom(() => {
        el.after(clone2);
        skipDuringClone(() => initTree(clone2))();
      });
      el._x_currentIfEl = clone2;
      el._x_undoIf = () => {
        walk(clone2, (node) => {
          if (!!node._x_effects) {
            node._x_effects.forEach(dequeueJob);
          }
        });
        clone2.remove();
        delete el._x_currentIfEl;
      };
      return clone2;
    };
    let hide = () => {
      if (!el._x_undoIf)
        return;
      el._x_undoIf();
      delete el._x_undoIf;
    };
    effect3(() => evaluate2((value) => {
      value ? show() : hide();
    }));
    cleanup2(() => el._x_undoIf && el._x_undoIf());
  });
  directive("id", (el, { expression }, { evaluate: evaluate2 }) => {
    let names = evaluate2(expression);
    names.forEach((name) => setIdRoot(el, name));
  });
  interceptClone((from, to) => {
    if (from._x_ids) {
      to._x_ids = from._x_ids;
    }
  });
  mapAttributes(startingWith("@", into(prefix("on:"))));
  directive("on", skipDuringClone((el, { value, modifiers, expression }, { cleanup: cleanup2 }) => {
    let evaluate2 = expression ? evaluateLater(el, expression) : () => {
    };
    if (el.tagName.toLowerCase() === "template") {
      if (!el._x_forwardEvents)
        el._x_forwardEvents = [];
      if (!el._x_forwardEvents.includes(value))
        el._x_forwardEvents.push(value);
    }
    let removeListener2 = on(el, value, modifiers, (e) => {
      evaluate2(() => {
      }, { scope: { "$event": e }, params: [e] });
    });
    cleanup2(() => removeListener2());
  }));
  warnMissingPluginDirective("Collapse", "collapse", "collapse");
  warnMissingPluginDirective("Intersect", "intersect", "intersect");
  warnMissingPluginDirective("Focus", "trap", "focus");
  warnMissingPluginDirective("Mask", "mask", "mask");
  function warnMissingPluginDirective(name, directiveName, slug) {
    directive(directiveName, (el) => warn(`You can't use [x-${directiveName}] without first installing the "${name}" plugin here: https://alpinejs.dev/plugins/${slug}`, el));
  }
  alpine_default.setEvaluator(normalEvaluator);
  alpine_default.setReactivityEngine({ reactive: reactive2, effect: effect2, release: stop, raw: toRaw });
  var src_default = alpine_default;
  var module_default = src_default;

  // node_modules/.pnpm/centrifuge@5.2.2/node_modules/centrifuge/build/index.mjs
  function __awaiter(thisArg, _arguments, P, generator) {
    function adopt(value) {
      return value instanceof P ? value : new P(function(resolve) {
        resolve(value);
      });
    }
    return new (P || (P = Promise))(function(resolve, reject) {
      function fulfilled(value) {
        try {
          step(generator.next(value));
        } catch (e) {
          reject(e);
        }
      }
      function rejected(value) {
        try {
          step(generator["throw"](value));
        } catch (e) {
          reject(e);
        }
      }
      function step(result) {
        result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected);
      }
      step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
  }
  function getDefaultExportFromCjs(x) {
    return x && x.__esModule && Object.prototype.hasOwnProperty.call(x, "default") ? x["default"] : x;
  }
  var events = { exports: {} };
  var R = typeof Reflect === "object" ? Reflect : null;
  var ReflectApply = R && typeof R.apply === "function" ? R.apply : function ReflectApply2(target, receiver, args) {
    return Function.prototype.apply.call(target, receiver, args);
  };
  var ReflectOwnKeys;
  if (R && typeof R.ownKeys === "function") {
    ReflectOwnKeys = R.ownKeys;
  } else if (Object.getOwnPropertySymbols) {
    ReflectOwnKeys = function ReflectOwnKeys2(target) {
      return Object.getOwnPropertyNames(target).concat(Object.getOwnPropertySymbols(target));
    };
  } else {
    ReflectOwnKeys = function ReflectOwnKeys2(target) {
      return Object.getOwnPropertyNames(target);
    };
  }
  function ProcessEmitWarning(warning) {
    if (console && console.warn) console.warn(warning);
  }
  var NumberIsNaN = Number.isNaN || function NumberIsNaN2(value) {
    return value !== value;
  };
  function EventEmitter() {
    EventEmitter.init.call(this);
  }
  events.exports = EventEmitter;
  events.exports.once = once3;
  EventEmitter.EventEmitter = EventEmitter;
  EventEmitter.prototype._events = void 0;
  EventEmitter.prototype._eventsCount = 0;
  EventEmitter.prototype._maxListeners = void 0;
  var defaultMaxListeners = 10;
  function checkListener(listener) {
    if (typeof listener !== "function") {
      throw new TypeError('The "listener" argument must be of type Function. Received type ' + typeof listener);
    }
  }
  Object.defineProperty(EventEmitter, "defaultMaxListeners", {
    enumerable: true,
    get: function() {
      return defaultMaxListeners;
    },
    set: function(arg) {
      if (typeof arg !== "number" || arg < 0 || NumberIsNaN(arg)) {
        throw new RangeError('The value of "defaultMaxListeners" is out of range. It must be a non-negative number. Received ' + arg + ".");
      }
      defaultMaxListeners = arg;
    }
  });
  EventEmitter.init = function() {
    if (this._events === void 0 || this._events === Object.getPrototypeOf(this)._events) {
      this._events = /* @__PURE__ */ Object.create(null);
      this._eventsCount = 0;
    }
    this._maxListeners = this._maxListeners || void 0;
  };
  EventEmitter.prototype.setMaxListeners = function setMaxListeners(n) {
    if (typeof n !== "number" || n < 0 || NumberIsNaN(n)) {
      throw new RangeError('The value of "n" is out of range. It must be a non-negative number. Received ' + n + ".");
    }
    this._maxListeners = n;
    return this;
  };
  function _getMaxListeners(that) {
    if (that._maxListeners === void 0)
      return EventEmitter.defaultMaxListeners;
    return that._maxListeners;
  }
  EventEmitter.prototype.getMaxListeners = function getMaxListeners() {
    return _getMaxListeners(this);
  };
  EventEmitter.prototype.emit = function emit(type) {
    var args = [];
    for (var i = 1; i < arguments.length; i++) args.push(arguments[i]);
    var doError = type === "error";
    var events2 = this._events;
    if (events2 !== void 0)
      doError = doError && events2.error === void 0;
    else if (!doError)
      return false;
    if (doError) {
      var er;
      if (args.length > 0)
        er = args[0];
      if (er instanceof Error) {
        throw er;
      }
      var err = new Error("Unhandled error." + (er ? " (" + er.message + ")" : ""));
      err.context = er;
      throw err;
    }
    var handler4 = events2[type];
    if (handler4 === void 0)
      return false;
    if (typeof handler4 === "function") {
      ReflectApply(handler4, this, args);
    } else {
      var len = handler4.length;
      var listeners2 = arrayClone(handler4, len);
      for (var i = 0; i < len; ++i)
        ReflectApply(listeners2[i], this, args);
    }
    return true;
  };
  function _addListener(target, type, listener, prepend) {
    var m;
    var events2;
    var existing;
    checkListener(listener);
    events2 = target._events;
    if (events2 === void 0) {
      events2 = target._events = /* @__PURE__ */ Object.create(null);
      target._eventsCount = 0;
    } else {
      if (events2.newListener !== void 0) {
        target.emit(
          "newListener",
          type,
          listener.listener ? listener.listener : listener
        );
        events2 = target._events;
      }
      existing = events2[type];
    }
    if (existing === void 0) {
      existing = events2[type] = listener;
      ++target._eventsCount;
    } else {
      if (typeof existing === "function") {
        existing = events2[type] = prepend ? [listener, existing] : [existing, listener];
      } else if (prepend) {
        existing.unshift(listener);
      } else {
        existing.push(listener);
      }
      m = _getMaxListeners(target);
      if (m > 0 && existing.length > m && !existing.warned) {
        existing.warned = true;
        var w = new Error("Possible EventEmitter memory leak detected. " + existing.length + " " + String(type) + " listeners added. Use emitter.setMaxListeners() to increase limit");
        w.name = "MaxListenersExceededWarning";
        w.emitter = target;
        w.type = type;
        w.count = existing.length;
        ProcessEmitWarning(w);
      }
    }
    return target;
  }
  EventEmitter.prototype.addListener = function addListener(type, listener) {
    return _addListener(this, type, listener, false);
  };
  EventEmitter.prototype.on = EventEmitter.prototype.addListener;
  EventEmitter.prototype.prependListener = function prependListener(type, listener) {
    return _addListener(this, type, listener, true);
  };
  function onceWrapper() {
    if (!this.fired) {
      this.target.removeListener(this.type, this.wrapFn);
      this.fired = true;
      if (arguments.length === 0)
        return this.listener.call(this.target);
      return this.listener.apply(this.target, arguments);
    }
  }
  function _onceWrap(target, type, listener) {
    var state = { fired: false, wrapFn: void 0, target, type, listener };
    var wrapped = onceWrapper.bind(state);
    wrapped.listener = listener;
    state.wrapFn = wrapped;
    return wrapped;
  }
  EventEmitter.prototype.once = function once2(type, listener) {
    checkListener(listener);
    this.on(type, _onceWrap(this, type, listener));
    return this;
  };
  EventEmitter.prototype.prependOnceListener = function prependOnceListener(type, listener) {
    checkListener(listener);
    this.prependListener(type, _onceWrap(this, type, listener));
    return this;
  };
  EventEmitter.prototype.removeListener = function removeListener(type, listener) {
    var list, events2, position, i, originalListener;
    checkListener(listener);
    events2 = this._events;
    if (events2 === void 0)
      return this;
    list = events2[type];
    if (list === void 0)
      return this;
    if (list === listener || list.listener === listener) {
      if (--this._eventsCount === 0)
        this._events = /* @__PURE__ */ Object.create(null);
      else {
        delete events2[type];
        if (events2.removeListener)
          this.emit("removeListener", type, list.listener || listener);
      }
    } else if (typeof list !== "function") {
      position = -1;
      for (i = list.length - 1; i >= 0; i--) {
        if (list[i] === listener || list[i].listener === listener) {
          originalListener = list[i].listener;
          position = i;
          break;
        }
      }
      if (position < 0)
        return this;
      if (position === 0)
        list.shift();
      else {
        spliceOne(list, position);
      }
      if (list.length === 1)
        events2[type] = list[0];
      if (events2.removeListener !== void 0)
        this.emit("removeListener", type, originalListener || listener);
    }
    return this;
  };
  EventEmitter.prototype.off = EventEmitter.prototype.removeListener;
  EventEmitter.prototype.removeAllListeners = function removeAllListeners(type) {
    var listeners2, events2, i;
    events2 = this._events;
    if (events2 === void 0)
      return this;
    if (events2.removeListener === void 0) {
      if (arguments.length === 0) {
        this._events = /* @__PURE__ */ Object.create(null);
        this._eventsCount = 0;
      } else if (events2[type] !== void 0) {
        if (--this._eventsCount === 0)
          this._events = /* @__PURE__ */ Object.create(null);
        else
          delete events2[type];
      }
      return this;
    }
    if (arguments.length === 0) {
      var keys = Object.keys(events2);
      var key;
      for (i = 0; i < keys.length; ++i) {
        key = keys[i];
        if (key === "removeListener") continue;
        this.removeAllListeners(key);
      }
      this.removeAllListeners("removeListener");
      this._events = /* @__PURE__ */ Object.create(null);
      this._eventsCount = 0;
      return this;
    }
    listeners2 = events2[type];
    if (typeof listeners2 === "function") {
      this.removeListener(type, listeners2);
    } else if (listeners2 !== void 0) {
      for (i = listeners2.length - 1; i >= 0; i--) {
        this.removeListener(type, listeners2[i]);
      }
    }
    return this;
  };
  function _listeners(target, type, unwrap) {
    var events2 = target._events;
    if (events2 === void 0)
      return [];
    var evlistener = events2[type];
    if (evlistener === void 0)
      return [];
    if (typeof evlistener === "function")
      return unwrap ? [evlistener.listener || evlistener] : [evlistener];
    return unwrap ? unwrapListeners(evlistener) : arrayClone(evlistener, evlistener.length);
  }
  EventEmitter.prototype.listeners = function listeners(type) {
    return _listeners(this, type, true);
  };
  EventEmitter.prototype.rawListeners = function rawListeners(type) {
    return _listeners(this, type, false);
  };
  EventEmitter.listenerCount = function(emitter, type) {
    if (typeof emitter.listenerCount === "function") {
      return emitter.listenerCount(type);
    } else {
      return listenerCount.call(emitter, type);
    }
  };
  EventEmitter.prototype.listenerCount = listenerCount;
  function listenerCount(type) {
    var events2 = this._events;
    if (events2 !== void 0) {
      var evlistener = events2[type];
      if (typeof evlistener === "function") {
        return 1;
      } else if (evlistener !== void 0) {
        return evlistener.length;
      }
    }
    return 0;
  }
  EventEmitter.prototype.eventNames = function eventNames() {
    return this._eventsCount > 0 ? ReflectOwnKeys(this._events) : [];
  };
  function arrayClone(arr, n) {
    var copy = new Array(n);
    for (var i = 0; i < n; ++i)
      copy[i] = arr[i];
    return copy;
  }
  function spliceOne(list, index) {
    for (; index + 1 < list.length; index++)
      list[index] = list[index + 1];
    list.pop();
  }
  function unwrapListeners(arr) {
    var ret = new Array(arr.length);
    for (var i = 0; i < ret.length; ++i) {
      ret[i] = arr[i].listener || arr[i];
    }
    return ret;
  }
  function once3(emitter, name) {
    return new Promise(function(resolve, reject) {
      function errorListener(err) {
        emitter.removeListener(name, resolver);
        reject(err);
      }
      function resolver() {
        if (typeof emitter.removeListener === "function") {
          emitter.removeListener("error", errorListener);
        }
        resolve([].slice.call(arguments));
      }
      eventTargetAgnosticAddListener(emitter, name, resolver, { once: true });
      if (name !== "error") {
        addErrorHandlerIfEventEmitter(emitter, errorListener, { once: true });
      }
    });
  }
  function addErrorHandlerIfEventEmitter(emitter, handler4, flags) {
    if (typeof emitter.on === "function") {
      eventTargetAgnosticAddListener(emitter, "error", handler4, flags);
    }
  }
  function eventTargetAgnosticAddListener(emitter, name, listener, flags) {
    if (typeof emitter.on === "function") {
      if (flags.once) {
        emitter.once(name, listener);
      } else {
        emitter.on(name, listener);
      }
    } else if (typeof emitter.addEventListener === "function") {
      emitter.addEventListener(name, function wrapListener(arg) {
        if (flags.once) {
          emitter.removeEventListener(name, wrapListener);
        }
        listener(arg);
      });
    } else {
      throw new TypeError('The "emitter" argument must be of type EventEmitter. Received type ' + typeof emitter);
    }
  }
  var eventsExports = events.exports;
  var EventEmitter$1 = /* @__PURE__ */ getDefaultExportFromCjs(eventsExports);
  var errorCodes;
  (function(errorCodes2) {
    errorCodes2[errorCodes2["timeout"] = 1] = "timeout";
    errorCodes2[errorCodes2["transportClosed"] = 2] = "transportClosed";
    errorCodes2[errorCodes2["clientDisconnected"] = 3] = "clientDisconnected";
    errorCodes2[errorCodes2["clientClosed"] = 4] = "clientClosed";
    errorCodes2[errorCodes2["clientConnectToken"] = 5] = "clientConnectToken";
    errorCodes2[errorCodes2["clientRefreshToken"] = 6] = "clientRefreshToken";
    errorCodes2[errorCodes2["subscriptionUnsubscribed"] = 7] = "subscriptionUnsubscribed";
    errorCodes2[errorCodes2["subscriptionSubscribeToken"] = 8] = "subscriptionSubscribeToken";
    errorCodes2[errorCodes2["subscriptionRefreshToken"] = 9] = "subscriptionRefreshToken";
    errorCodes2[errorCodes2["transportWriteError"] = 10] = "transportWriteError";
    errorCodes2[errorCodes2["connectionClosed"] = 11] = "connectionClosed";
    errorCodes2[errorCodes2["badConfiguration"] = 12] = "badConfiguration";
  })(errorCodes || (errorCodes = {}));
  var connectingCodes;
  (function(connectingCodes2) {
    connectingCodes2[connectingCodes2["connectCalled"] = 0] = "connectCalled";
    connectingCodes2[connectingCodes2["transportClosed"] = 1] = "transportClosed";
    connectingCodes2[connectingCodes2["noPing"] = 2] = "noPing";
    connectingCodes2[connectingCodes2["subscribeTimeout"] = 3] = "subscribeTimeout";
    connectingCodes2[connectingCodes2["unsubscribeError"] = 4] = "unsubscribeError";
  })(connectingCodes || (connectingCodes = {}));
  var disconnectedCodes;
  (function(disconnectedCodes2) {
    disconnectedCodes2[disconnectedCodes2["disconnectCalled"] = 0] = "disconnectCalled";
    disconnectedCodes2[disconnectedCodes2["unauthorized"] = 1] = "unauthorized";
    disconnectedCodes2[disconnectedCodes2["badProtocol"] = 2] = "badProtocol";
    disconnectedCodes2[disconnectedCodes2["messageSizeLimit"] = 3] = "messageSizeLimit";
  })(disconnectedCodes || (disconnectedCodes = {}));
  var subscribingCodes;
  (function(subscribingCodes2) {
    subscribingCodes2[subscribingCodes2["subscribeCalled"] = 0] = "subscribeCalled";
    subscribingCodes2[subscribingCodes2["transportClosed"] = 1] = "transportClosed";
  })(subscribingCodes || (subscribingCodes = {}));
  var unsubscribedCodes;
  (function(unsubscribedCodes2) {
    unsubscribedCodes2[unsubscribedCodes2["unsubscribeCalled"] = 0] = "unsubscribeCalled";
    unsubscribedCodes2[unsubscribedCodes2["unauthorized"] = 1] = "unauthorized";
    unsubscribedCodes2[unsubscribedCodes2["clientClosed"] = 2] = "clientClosed";
  })(unsubscribedCodes || (unsubscribedCodes = {}));
  var State;
  (function(State2) {
    State2["Disconnected"] = "disconnected";
    State2["Connecting"] = "connecting";
    State2["Connected"] = "connected";
  })(State || (State = {}));
  var SubscriptionState;
  (function(SubscriptionState2) {
    SubscriptionState2["Unsubscribed"] = "unsubscribed";
    SubscriptionState2["Subscribing"] = "subscribing";
    SubscriptionState2["Subscribed"] = "subscribed";
  })(SubscriptionState || (SubscriptionState = {}));
  function startsWith(value, prefix2) {
    return value.lastIndexOf(prefix2, 0) === 0;
  }
  function isFunction(value) {
    if (value === void 0 || value === null) {
      return false;
    }
    return typeof value === "function";
  }
  function log(level, args) {
    if (globalThis.console) {
      const logger = globalThis.console[level];
      if (isFunction(logger)) {
        logger.apply(globalThis.console, args);
      }
    }
  }
  function randomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1) + min);
  }
  function backoff(step, min, max) {
    if (step > 31) {
      step = 31;
    }
    const interval = randomInt(0, Math.min(max, min * Math.pow(2, step)));
    return Math.min(max, min + interval);
  }
  function errorExists(data2) {
    return "error" in data2 && data2.error !== null;
  }
  function ttlMilliseconds(ttl) {
    return Math.min(ttl * 1e3, 2147483647);
  }
  var Subscription = class extends EventEmitter$1 {
    /** Subscription constructor should not be used directly, create subscriptions using Client method. */
    constructor(centrifuge, channel, options) {
      super();
      this._resubscribeTimeout = null;
      this._refreshTimeout = null;
      this.channel = channel;
      this.state = SubscriptionState.Unsubscribed;
      this._centrifuge = centrifuge;
      this._token = "";
      this._getToken = null;
      this._data = null;
      this._getData = null;
      this._recover = false;
      this._offset = null;
      this._epoch = null;
      this._recoverable = false;
      this._positioned = false;
      this._joinLeave = false;
      this._minResubscribeDelay = 500;
      this._maxResubscribeDelay = 2e4;
      this._resubscribeTimeout = null;
      this._resubscribeAttempts = 0;
      this._promises = {};
      this._promiseId = 0;
      this._inflight = false;
      this._refreshTimeout = null;
      this._delta = "";
      this._delta_negotiated = false;
      this._prevValue = null;
      this._unsubPromise = Promise.resolve();
      this._setOptions(options);
      if (this._centrifuge._debugEnabled) {
        this.on("state", (ctx) => {
          this._centrifuge._debug("subscription state", channel, ctx.oldState, "->", ctx.newState);
        });
        this.on("error", (ctx) => {
          this._centrifuge._debug("subscription error", channel, ctx);
        });
      } else {
        this.on("error", function() {
          Function.prototype();
        });
      }
    }
    /** ready returns a Promise which resolves upon subscription goes to Subscribed
     * state and rejects in case of subscription goes to Unsubscribed state.
     * Optional timeout can be passed.*/
    ready(timeout) {
      if (this.state === SubscriptionState.Unsubscribed) {
        return Promise.reject({ code: errorCodes.subscriptionUnsubscribed, message: this.state });
      }
      if (this.state === SubscriptionState.Subscribed) {
        return Promise.resolve();
      }
      return new Promise((res, rej) => {
        const ctx = {
          resolve: res,
          reject: rej
        };
        if (timeout) {
          ctx.timeout = setTimeout(function() {
            rej({ code: errorCodes.timeout, message: "timeout" });
          }, timeout);
        }
        this._promises[this._nextPromiseId()] = ctx;
      });
    }
    /** subscribe to a channel.*/
    subscribe() {
      if (this._isSubscribed()) {
        return;
      }
      this._resubscribeAttempts = 0;
      this._setSubscribing(subscribingCodes.subscribeCalled, "subscribe called");
    }
    /** unsubscribe from a channel, keeping position state.*/
    unsubscribe() {
      this._unsubPromise = this._setUnsubscribed(unsubscribedCodes.unsubscribeCalled, "unsubscribe called", true);
    }
    /** publish data to a channel.*/
    publish(data2) {
      const self = this;
      return this._methodCall().then(function() {
        return self._centrifuge.publish(self.channel, data2);
      });
    }
    /** get online presence for a channel.*/
    presence() {
      const self = this;
      return this._methodCall().then(function() {
        return self._centrifuge.presence(self.channel);
      });
    }
    /** presence stats for a channel (num clients and unique users).*/
    presenceStats() {
      const self = this;
      return this._methodCall().then(function() {
        return self._centrifuge.presenceStats(self.channel);
      });
    }
    /** history for a channel. By default it does not return publications (only current
     *  StreamPosition data) – provide an explicit limit > 0 to load publications.*/
    history(opts) {
      const self = this;
      return this._methodCall().then(function() {
        return self._centrifuge.history(self.channel, opts);
      });
    }
    _methodCall() {
      if (this._isSubscribed()) {
        return Promise.resolve();
      }
      if (this._isUnsubscribed()) {
        return Promise.reject({ code: errorCodes.subscriptionUnsubscribed, message: this.state });
      }
      return new Promise((res, rej) => {
        const timeout = setTimeout(function() {
          rej({ code: errorCodes.timeout, message: "timeout" });
        }, this._centrifuge._config.timeout);
        this._promises[this._nextPromiseId()] = {
          timeout,
          resolve: res,
          reject: rej
        };
      });
    }
    _nextPromiseId() {
      return ++this._promiseId;
    }
    _needRecover() {
      return this._recover === true;
    }
    _isUnsubscribed() {
      return this.state === SubscriptionState.Unsubscribed;
    }
    _isSubscribing() {
      return this.state === SubscriptionState.Subscribing;
    }
    _isSubscribed() {
      return this.state === SubscriptionState.Subscribed;
    }
    _setState(newState) {
      if (this.state !== newState) {
        const oldState = this.state;
        this.state = newState;
        this.emit("state", { newState, oldState, channel: this.channel });
        return true;
      }
      return false;
    }
    _usesToken() {
      return this._token !== "" || this._getToken !== null;
    }
    _clearSubscribingState() {
      this._resubscribeAttempts = 0;
      this._clearResubscribeTimeout();
    }
    _clearSubscribedState() {
      this._clearRefreshTimeout();
    }
    _setSubscribed(result) {
      if (!this._isSubscribing()) {
        return;
      }
      this._clearSubscribingState();
      if (result.recoverable) {
        this._recover = true;
        this._offset = result.offset || 0;
        this._epoch = result.epoch || "";
      }
      if (result.delta) {
        this._delta_negotiated = true;
      } else {
        this._delta_negotiated = false;
      }
      this._setState(SubscriptionState.Subscribed);
      const ctx = this._centrifuge._getSubscribeContext(this.channel, result);
      this.emit("subscribed", ctx);
      this._resolvePromises();
      const pubs = result.publications;
      if (pubs && pubs.length > 0) {
        for (const i in pubs) {
          if (!pubs.hasOwnProperty(i)) {
            continue;
          }
          this._handlePublication(pubs[i]);
        }
      }
      if (result.expires === true) {
        this._refreshTimeout = setTimeout(() => this._refresh(), ttlMilliseconds(result.ttl));
      }
    }
    _setSubscribing(code, reason) {
      return __awaiter(this, void 0, void 0, function* () {
        if (this._isSubscribing()) {
          return;
        }
        if (this._isSubscribed()) {
          this._clearSubscribedState();
        }
        if (this._setState(SubscriptionState.Subscribing)) {
          this.emit("subscribing", { channel: this.channel, code, reason });
        }
        if (this._centrifuge._transport && this._centrifuge._transport.emulation()) {
          yield this._unsubPromise;
        }
        if (!this._isSubscribing()) {
          return;
        }
        this._subscribe();
      });
    }
    _subscribe() {
      this._centrifuge._debug("subscribing on", this.channel);
      if (!this._centrifuge._transportIsOpen) {
        this._centrifuge._debug("delay subscribe on", this.channel, "till connected");
        return null;
      }
      const self = this;
      const getDataCtx = {
        channel: self.channel
      };
      if (!this._usesToken() || this._token) {
        if (self._getData) {
          self._getData(getDataCtx).then(function(data2) {
            if (!self._isSubscribing()) {
              return;
            }
            self._data = data2;
            self._sendSubscribe(self._token);
          });
          return null;
        } else {
          return self._sendSubscribe(self._token);
        }
      }
      this._getSubscriptionToken().then(function(token) {
        if (!self._isSubscribing()) {
          return;
        }
        if (!token) {
          self._failUnauthorized();
          return;
        }
        self._token = token;
        if (self._getData) {
          self._getData(getDataCtx).then(function(data2) {
            if (!self._isSubscribing()) {
              return;
            }
            self._data = data2;
            self._sendSubscribe(token);
          });
        } else {
          self._sendSubscribe(token);
        }
      }).catch(function(e) {
        if (!self._isSubscribing()) {
          return;
        }
        if (e instanceof UnauthorizedError) {
          self._failUnauthorized();
          return;
        }
        self.emit("error", {
          type: "subscribeToken",
          channel: self.channel,
          error: {
            code: errorCodes.subscriptionSubscribeToken,
            message: e !== void 0 ? e.toString() : ""
          }
        });
        self._scheduleResubscribe();
      });
      return null;
    }
    _sendSubscribe(token) {
      if (!this._centrifuge._transportIsOpen) {
        return null;
      }
      const channel = this.channel;
      const req = {
        channel
      };
      if (token) {
        req.token = token;
      }
      if (this._data) {
        req.data = this._data;
      }
      if (this._positioned) {
        req.positioned = true;
      }
      if (this._recoverable) {
        req.recoverable = true;
      }
      if (this._joinLeave) {
        req.join_leave = true;
      }
      if (this._needRecover()) {
        req.recover = true;
        const offset = this._getOffset();
        if (offset) {
          req.offset = offset;
        }
        const epoch = this._getEpoch();
        if (epoch) {
          req.epoch = epoch;
        }
      }
      if (this._delta) {
        req.delta = this._delta;
      }
      const cmd = { subscribe: req };
      this._inflight = true;
      this._centrifuge._call(cmd).then((resolveCtx) => {
        this._inflight = false;
        const result = resolveCtx.reply.subscribe;
        this._handleSubscribeResponse(result);
        if (resolveCtx.next) {
          resolveCtx.next();
        }
      }, (rejectCtx) => {
        this._inflight = false;
        this._handleSubscribeError(rejectCtx.error);
        if (rejectCtx.next) {
          rejectCtx.next();
        }
      });
      return cmd;
    }
    _handleSubscribeError(error2) {
      if (!this._isSubscribing()) {
        return;
      }
      if (error2.code === errorCodes.timeout) {
        this._centrifuge._disconnect(connectingCodes.subscribeTimeout, "subscribe timeout", true);
        return;
      }
      this._subscribeError(error2);
    }
    _handleSubscribeResponse(result) {
      if (!this._isSubscribing()) {
        return;
      }
      this._setSubscribed(result);
    }
    _setUnsubscribed(code, reason, sendUnsubscribe) {
      if (this._isUnsubscribed()) {
        return Promise.resolve();
      }
      let promise = Promise.resolve();
      if (this._isSubscribed()) {
        if (sendUnsubscribe) {
          promise = this._centrifuge._unsubscribe(this);
        }
        this._clearSubscribedState();
      } else if (this._isSubscribing()) {
        if (this._inflight && sendUnsubscribe) {
          promise = this._centrifuge._unsubscribe(this);
        }
        this._clearSubscribingState();
      }
      if (this._setState(SubscriptionState.Unsubscribed)) {
        this.emit("unsubscribed", { channel: this.channel, code, reason });
      }
      this._rejectPromises({ code: errorCodes.subscriptionUnsubscribed, message: this.state });
      return promise;
    }
    _handlePublication(pub) {
      if (this._delta && this._delta_negotiated) {
        const { newData, newPrevValue } = this._centrifuge._codec.applyDeltaIfNeeded(pub, this._prevValue);
        pub.data = newData;
        this._prevValue = newPrevValue;
      }
      const ctx = this._centrifuge._getPublicationContext(this.channel, pub);
      this.emit("publication", ctx);
      if (pub.offset) {
        this._offset = pub.offset;
      }
    }
    _handleJoin(join) {
      const info = this._centrifuge._getJoinLeaveContext(join.info);
      this.emit("join", { channel: this.channel, info });
    }
    _handleLeave(leave) {
      const info = this._centrifuge._getJoinLeaveContext(leave.info);
      this.emit("leave", { channel: this.channel, info });
    }
    _resolvePromises() {
      for (const id in this._promises) {
        if (!this._promises.hasOwnProperty(id)) {
          continue;
        }
        if (this._promises[id].timeout) {
          clearTimeout(this._promises[id].timeout);
        }
        this._promises[id].resolve();
        delete this._promises[id];
      }
    }
    _rejectPromises(err) {
      for (const id in this._promises) {
        if (!this._promises.hasOwnProperty(id)) {
          continue;
        }
        if (this._promises[id].timeout) {
          clearTimeout(this._promises[id].timeout);
        }
        this._promises[id].reject(err);
        delete this._promises[id];
      }
    }
    _scheduleResubscribe() {
      const self = this;
      const delay = this._getResubscribeDelay();
      this._resubscribeTimeout = setTimeout(function() {
        if (self._isSubscribing()) {
          self._subscribe();
        }
      }, delay);
    }
    _subscribeError(err) {
      if (!this._isSubscribing()) {
        return;
      }
      if (err.code < 100 || err.code === 109 || err.temporary === true) {
        if (err.code === 109) {
          this._token = "";
        }
        const errContext = {
          channel: this.channel,
          type: "subscribe",
          error: err
        };
        if (this._centrifuge.state === State.Connected) {
          this.emit("error", errContext);
        }
        this._scheduleResubscribe();
      } else {
        this._setUnsubscribed(err.code, err.message, false);
      }
    }
    _getResubscribeDelay() {
      const delay = backoff(this._resubscribeAttempts, this._minResubscribeDelay, this._maxResubscribeDelay);
      this._resubscribeAttempts++;
      return delay;
    }
    _setOptions(options) {
      if (!options) {
        return;
      }
      if (options.since) {
        this._offset = options.since.offset;
        this._epoch = options.since.epoch;
        this._recover = true;
      }
      if (options.data) {
        this._data = options.data;
      }
      if (options.getData) {
        this._getData = options.getData;
      }
      if (options.minResubscribeDelay !== void 0) {
        this._minResubscribeDelay = options.minResubscribeDelay;
      }
      if (options.maxResubscribeDelay !== void 0) {
        this._maxResubscribeDelay = options.maxResubscribeDelay;
      }
      if (options.token) {
        this._token = options.token;
      }
      if (options.getToken) {
        this._getToken = options.getToken;
      }
      if (options.positioned === true) {
        this._positioned = true;
      }
      if (options.recoverable === true) {
        this._recoverable = true;
      }
      if (options.joinLeave === true) {
        this._joinLeave = true;
      }
      if (options.delta) {
        if (options.delta !== "fossil") {
          throw new Error("unsupported delta format");
        }
        this._delta = options.delta;
      }
    }
    _getOffset() {
      const offset = this._offset;
      if (offset !== null) {
        return offset;
      }
      return 0;
    }
    _getEpoch() {
      const epoch = this._epoch;
      if (epoch !== null) {
        return epoch;
      }
      return "";
    }
    _clearRefreshTimeout() {
      if (this._refreshTimeout !== null) {
        clearTimeout(this._refreshTimeout);
        this._refreshTimeout = null;
      }
    }
    _clearResubscribeTimeout() {
      if (this._resubscribeTimeout !== null) {
        clearTimeout(this._resubscribeTimeout);
        this._resubscribeTimeout = null;
      }
    }
    _getSubscriptionToken() {
      this._centrifuge._debug("get subscription token for channel", this.channel);
      const ctx = {
        channel: this.channel
      };
      const getToken = this._getToken;
      if (getToken === null) {
        this.emit("error", {
          type: "configuration",
          channel: this.channel,
          error: {
            code: errorCodes.badConfiguration,
            message: "provide a function to get channel subscription token"
          }
        });
        throw new UnauthorizedError("");
      }
      return getToken(ctx);
    }
    _refresh() {
      this._clearRefreshTimeout();
      const self = this;
      this._getSubscriptionToken().then(function(token) {
        if (!self._isSubscribed()) {
          return;
        }
        if (!token) {
          self._failUnauthorized();
          return;
        }
        self._token = token;
        const req = {
          channel: self.channel,
          token
        };
        const msg = {
          "sub_refresh": req
        };
        self._centrifuge._call(msg).then((resolveCtx) => {
          const result = resolveCtx.reply.sub_refresh;
          self._refreshResponse(result);
          if (resolveCtx.next) {
            resolveCtx.next();
          }
        }, (rejectCtx) => {
          self._refreshError(rejectCtx.error);
          if (rejectCtx.next) {
            rejectCtx.next();
          }
        });
      }).catch(function(e) {
        if (e instanceof UnauthorizedError) {
          self._failUnauthorized();
          return;
        }
        self.emit("error", {
          type: "refreshToken",
          channel: self.channel,
          error: {
            code: errorCodes.subscriptionRefreshToken,
            message: e !== void 0 ? e.toString() : ""
          }
        });
        self._refreshTimeout = setTimeout(() => self._refresh(), self._getRefreshRetryDelay());
      });
    }
    _refreshResponse(result) {
      if (!this._isSubscribed()) {
        return;
      }
      this._centrifuge._debug("subscription token refreshed, channel", this.channel);
      this._clearRefreshTimeout();
      if (result.expires === true) {
        this._refreshTimeout = setTimeout(() => this._refresh(), ttlMilliseconds(result.ttl));
      }
    }
    _refreshError(err) {
      if (!this._isSubscribed()) {
        return;
      }
      if (err.code < 100 || err.temporary === true) {
        this.emit("error", {
          type: "refresh",
          channel: this.channel,
          error: err
        });
        this._refreshTimeout = setTimeout(() => this._refresh(), this._getRefreshRetryDelay());
      } else {
        this._setUnsubscribed(err.code, err.message, true);
      }
    }
    _getRefreshRetryDelay() {
      return backoff(0, 1e4, 2e4);
    }
    _failUnauthorized() {
      this._setUnsubscribed(unsubscribedCodes.unauthorized, "unauthorized", true);
    }
  };
  var SockjsTransport = class {
    constructor(endpoint, options) {
      this.endpoint = endpoint;
      this.options = options;
      this._transport = null;
    }
    name() {
      return "sockjs";
    }
    subName() {
      return "sockjs-" + this._transport.transport;
    }
    emulation() {
      return false;
    }
    supported() {
      return this.options.sockjs !== null;
    }
    initialize(_protocol, callbacks) {
      this._transport = new this.options.sockjs(this.endpoint, null, this.options.sockjsOptions);
      this._transport.onopen = () => {
        callbacks.onOpen();
      };
      this._transport.onerror = (e) => {
        callbacks.onError(e);
      };
      this._transport.onclose = (closeEvent) => {
        callbacks.onClose(closeEvent);
      };
      this._transport.onmessage = (event) => {
        callbacks.onMessage(event.data);
      };
    }
    close() {
      this._transport.close();
    }
    send(data2) {
      this._transport.send(data2);
    }
  };
  var WebsocketTransport = class {
    constructor(endpoint, options) {
      this.endpoint = endpoint;
      this.options = options;
      this._transport = null;
    }
    name() {
      return "websocket";
    }
    subName() {
      return "websocket";
    }
    emulation() {
      return false;
    }
    supported() {
      return this.options.websocket !== void 0 && this.options.websocket !== null;
    }
    initialize(protocol, callbacks) {
      let subProtocol = "";
      if (protocol === "protobuf") {
        subProtocol = "centrifuge-protobuf";
      }
      if (subProtocol !== "") {
        this._transport = new this.options.websocket(this.endpoint, subProtocol);
      } else {
        this._transport = new this.options.websocket(this.endpoint);
      }
      if (protocol === "protobuf") {
        this._transport.binaryType = "arraybuffer";
      }
      this._transport.onopen = () => {
        callbacks.onOpen();
      };
      this._transport.onerror = (e) => {
        callbacks.onError(e);
      };
      this._transport.onclose = (closeEvent) => {
        callbacks.onClose(closeEvent);
      };
      this._transport.onmessage = (event) => {
        callbacks.onMessage(event.data);
      };
    }
    close() {
      this._transport.close();
    }
    send(data2) {
      this._transport.send(data2);
    }
  };
  var HttpStreamTransport = class {
    constructor(endpoint, options) {
      this.endpoint = endpoint;
      this.options = options;
      this._abortController = null;
      this._utf8decoder = new TextDecoder();
      this._protocol = "json";
    }
    name() {
      return "http_stream";
    }
    subName() {
      return "http_stream";
    }
    emulation() {
      return true;
    }
    _handleErrors(response) {
      if (!response.ok)
        throw new Error(response.status);
      return response;
    }
    _fetchEventTarget(self, endpoint, options) {
      const eventTarget = new EventTarget();
      const fetchFunc = self.options.fetch;
      fetchFunc(endpoint, options).then(self._handleErrors).then((response) => {
        eventTarget.dispatchEvent(new Event("open"));
        let jsonStreamBuf = "";
        let jsonStreamPos = 0;
        let protoStreamBuf = new Uint8Array();
        const reader = response.body.getReader();
        return new self.options.readableStream({
          start(controller) {
            function pump() {
              return reader.read().then(({ done, value }) => {
                if (done) {
                  eventTarget.dispatchEvent(new Event("close"));
                  controller.close();
                  return;
                }
                try {
                  if (self._protocol === "json") {
                    jsonStreamBuf += self._utf8decoder.decode(value);
                    while (jsonStreamPos < jsonStreamBuf.length) {
                      if (jsonStreamBuf[jsonStreamPos] === "\n") {
                        const line = jsonStreamBuf.substring(0, jsonStreamPos);
                        eventTarget.dispatchEvent(new MessageEvent("message", { data: line }));
                        jsonStreamBuf = jsonStreamBuf.substring(jsonStreamPos + 1);
                        jsonStreamPos = 0;
                      } else {
                        ++jsonStreamPos;
                      }
                    }
                  } else {
                    const mergedArray = new Uint8Array(protoStreamBuf.length + value.length);
                    mergedArray.set(protoStreamBuf);
                    mergedArray.set(value, protoStreamBuf.length);
                    protoStreamBuf = mergedArray;
                    while (true) {
                      const result = self.options.decoder.decodeReply(protoStreamBuf);
                      if (result.ok) {
                        const data2 = protoStreamBuf.slice(0, result.pos);
                        eventTarget.dispatchEvent(new MessageEvent("message", { data: data2 }));
                        protoStreamBuf = protoStreamBuf.slice(result.pos);
                        continue;
                      }
                      break;
                    }
                  }
                } catch (error2) {
                  eventTarget.dispatchEvent(new Event("error", { detail: error2 }));
                  eventTarget.dispatchEvent(new Event("close"));
                  controller.close();
                  return;
                }
                pump();
              }).catch(function(e) {
                eventTarget.dispatchEvent(new Event("error", { detail: e }));
                eventTarget.dispatchEvent(new Event("close"));
                controller.close();
                return;
              });
            }
            return pump();
          }
        });
      }).catch((error2) => {
        eventTarget.dispatchEvent(new Event("error", { detail: error2 }));
        eventTarget.dispatchEvent(new Event("close"));
      });
      return eventTarget;
    }
    supported() {
      return this.options.fetch !== null && this.options.readableStream !== null && typeof TextDecoder !== "undefined" && typeof AbortController !== "undefined" && typeof EventTarget !== "undefined" && typeof Event !== "undefined" && typeof MessageEvent !== "undefined" && typeof Error !== "undefined";
    }
    initialize(protocol, callbacks, initialData) {
      this._protocol = protocol;
      this._abortController = new AbortController();
      let headers;
      let body;
      if (protocol === "json") {
        headers = {
          "Accept": "application/json",
          "Content-Type": "application/json"
        };
        body = initialData;
      } else {
        headers = {
          "Accept": "application/octet-stream",
          "Content-Type": "application/octet-stream"
        };
        body = initialData;
      }
      const fetchOptions = {
        method: "POST",
        headers,
        body,
        mode: "cors",
        credentials: "same-origin",
        cache: "no-cache",
        signal: this._abortController.signal
      };
      const eventTarget = this._fetchEventTarget(this, this.endpoint, fetchOptions);
      eventTarget.addEventListener("open", () => {
        callbacks.onOpen();
      });
      eventTarget.addEventListener("error", (e) => {
        this._abortController.abort();
        callbacks.onError(e);
      });
      eventTarget.addEventListener("close", () => {
        this._abortController.abort();
        callbacks.onClose({
          code: 4,
          reason: "connection closed"
        });
      });
      eventTarget.addEventListener("message", (e) => {
        callbacks.onMessage(e.data);
      });
    }
    close() {
      this._abortController.abort();
    }
    send(data2, session, node) {
      let headers;
      let body;
      const req = {
        session,
        node,
        data: data2
      };
      if (this._protocol === "json") {
        headers = {
          "Content-Type": "application/json"
        };
        body = JSON.stringify(req);
      } else {
        headers = {
          "Content-Type": "application/octet-stream"
        };
        body = this.options.encoder.encodeEmulationRequest(req);
      }
      const fetchFunc = this.options.fetch;
      const fetchOptions = {
        method: "POST",
        headers,
        body,
        mode: "cors",
        credentials: "same-origin",
        cache: "no-cache"
      };
      fetchFunc(this.options.emulationEndpoint, fetchOptions);
    }
  };
  var SseTransport = class {
    constructor(endpoint, options) {
      this.endpoint = endpoint;
      this.options = options;
      this._protocol = "json";
      this._transport = null;
      this._onClose = null;
    }
    name() {
      return "sse";
    }
    subName() {
      return "sse";
    }
    emulation() {
      return true;
    }
    supported() {
      return this.options.eventsource !== null && this.options.fetch !== null;
    }
    initialize(_protocol, callbacks, initialData) {
      let url;
      if (globalThis && globalThis.document && globalThis.document.baseURI) {
        url = new URL(this.endpoint, globalThis.document.baseURI);
      } else {
        url = new URL(this.endpoint);
      }
      url.searchParams.append("cf_connect", initialData);
      const eventsourceOptions = {};
      const eventSource = new this.options.eventsource(url.toString(), eventsourceOptions);
      this._transport = eventSource;
      const self = this;
      eventSource.onopen = function() {
        callbacks.onOpen();
      };
      eventSource.onerror = function(e) {
        eventSource.close();
        callbacks.onError(e);
        callbacks.onClose({
          code: 4,
          reason: "connection closed"
        });
      };
      eventSource.onmessage = function(e) {
        callbacks.onMessage(e.data);
      };
      self._onClose = function() {
        callbacks.onClose({
          code: 4,
          reason: "connection closed"
        });
      };
    }
    close() {
      this._transport.close();
      if (this._onClose !== null) {
        this._onClose();
      }
    }
    send(data2, session, node) {
      const req = {
        session,
        node,
        data: data2
      };
      const headers = {
        "Content-Type": "application/json"
      };
      const body = JSON.stringify(req);
      const fetchFunc = this.options.fetch;
      const fetchOptions = {
        method: "POST",
        headers,
        body,
        mode: "cors",
        credentials: "same-origin",
        cache: "no-cache"
      };
      fetchFunc(this.options.emulationEndpoint, fetchOptions);
    }
  };
  var WebtransportTransport = class {
    constructor(endpoint, options) {
      this.endpoint = endpoint;
      this.options = options;
      this._transport = null;
      this._stream = null;
      this._writer = null;
      this._utf8decoder = new TextDecoder();
      this._protocol = "json";
    }
    name() {
      return "webtransport";
    }
    subName() {
      return "webtransport";
    }
    emulation() {
      return false;
    }
    supported() {
      return this.options.webtransport !== void 0 && this.options.webtransport !== null;
    }
    initialize(protocol, callbacks) {
      return __awaiter(this, void 0, void 0, function* () {
        let url;
        if (globalThis && globalThis.document && globalThis.document.baseURI) {
          url = new URL(this.endpoint, globalThis.document.baseURI);
        } else {
          url = new URL(this.endpoint);
        }
        if (protocol === "protobuf") {
          url.searchParams.append("cf_protocol", "protobuf");
        }
        this._protocol = protocol;
        const eventTarget = new EventTarget();
        this._transport = new this.options.webtransport(url.toString());
        this._transport.closed.then(() => {
          callbacks.onClose({
            code: 4,
            reason: "connection closed"
          });
        }).catch(() => {
          callbacks.onClose({
            code: 4,
            reason: "connection closed"
          });
        });
        try {
          yield this._transport.ready;
        } catch (_a) {
          this.close();
          return;
        }
        let stream;
        try {
          stream = yield this._transport.createBidirectionalStream();
        } catch (_b) {
          this.close();
          return;
        }
        this._stream = stream;
        this._writer = this._stream.writable.getWriter();
        eventTarget.addEventListener("close", () => {
          callbacks.onClose({
            code: 4,
            reason: "connection closed"
          });
        });
        eventTarget.addEventListener("message", (e) => {
          callbacks.onMessage(e.data);
        });
        this._startReading(eventTarget);
        callbacks.onOpen();
      });
    }
    _startReading(eventTarget) {
      return __awaiter(this, void 0, void 0, function* () {
        const reader = this._stream.readable.getReader();
        let jsonStreamBuf = "";
        let jsonStreamPos = 0;
        let protoStreamBuf = new Uint8Array();
        try {
          while (true) {
            const { done, value } = yield reader.read();
            if (value.length > 0) {
              if (this._protocol === "json") {
                jsonStreamBuf += this._utf8decoder.decode(value);
                while (jsonStreamPos < jsonStreamBuf.length) {
                  if (jsonStreamBuf[jsonStreamPos] === "\n") {
                    const line = jsonStreamBuf.substring(0, jsonStreamPos);
                    eventTarget.dispatchEvent(new MessageEvent("message", { data: line }));
                    jsonStreamBuf = jsonStreamBuf.substring(jsonStreamPos + 1);
                    jsonStreamPos = 0;
                  } else {
                    ++jsonStreamPos;
                  }
                }
              } else {
                const mergedArray = new Uint8Array(protoStreamBuf.length + value.length);
                mergedArray.set(protoStreamBuf);
                mergedArray.set(value, protoStreamBuf.length);
                protoStreamBuf = mergedArray;
                while (true) {
                  const result = this.options.decoder.decodeReply(protoStreamBuf);
                  if (result.ok) {
                    const data2 = protoStreamBuf.slice(0, result.pos);
                    eventTarget.dispatchEvent(new MessageEvent("message", { data: data2 }));
                    protoStreamBuf = protoStreamBuf.slice(result.pos);
                    continue;
                  }
                  break;
                }
              }
            }
            if (done) {
              break;
            }
          }
        } catch (_a) {
          eventTarget.dispatchEvent(new Event("close"));
        }
      });
    }
    close() {
      return __awaiter(this, void 0, void 0, function* () {
        try {
          if (this._writer) {
            yield this._writer.close();
          }
          this._transport.close();
        } catch (e) {
        }
      });
    }
    send(data2) {
      return __awaiter(this, void 0, void 0, function* () {
        let binary;
        if (this._protocol === "json") {
          binary = new TextEncoder().encode(data2 + "\n");
        } else {
          binary = data2;
        }
        try {
          yield this._writer.write(binary);
        } catch (e) {
          this.close();
        }
      });
    }
  };
  var zValue = [
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    0,
    1,
    2,
    3,
    4,
    5,
    6,
    7,
    8,
    9,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    -1,
    10,
    11,
    12,
    13,
    14,
    15,
    16,
    17,
    18,
    19,
    20,
    21,
    22,
    23,
    24,
    25,
    26,
    27,
    28,
    29,
    30,
    31,
    32,
    33,
    34,
    35,
    -1,
    -1,
    -1,
    -1,
    36,
    -1,
    37,
    38,
    39,
    40,
    41,
    42,
    43,
    44,
    45,
    46,
    47,
    48,
    49,
    50,
    51,
    52,
    53,
    54,
    55,
    56,
    57,
    58,
    59,
    60,
    61,
    62,
    -1,
    -1,
    -1,
    63,
    -1
  ];
  var Reader = class {
    constructor(array) {
      this.a = array;
      this.pos = 0;
    }
    haveBytes() {
      return this.pos < this.a.length;
    }
    getByte() {
      const b = this.a[this.pos];
      this.pos++;
      if (this.pos > this.a.length)
        throw new RangeError("out of bounds");
      return b;
    }
    getChar() {
      return String.fromCharCode(this.getByte());
    }
    // Read base64-encoded unsigned integer.
    getInt() {
      let v = 0;
      let c;
      while (this.haveBytes() && (c = zValue[127 & this.getByte()]) >= 0) {
        v = (v << 6) + c;
      }
      this.pos--;
      return v >>> 0;
    }
  };
  var Writer = class {
    constructor() {
      this.a = [];
    }
    toByteArray(sourceType) {
      if (Array.isArray(sourceType)) {
        return this.a;
      }
      return new Uint8Array(this.a);
    }
    // Copy from array at start to end.
    putArray(a, start2, end) {
      for (let i = start2; i < end; i++)
        this.a.push(a[i]);
    }
  };
  function checksum(arr) {
    let sum0 = 0, sum1 = 0, sum2 = 0, sum3 = 0, z = 0, N = arr.length;
    while (N >= 16) {
      sum0 = sum0 + arr[z + 0] | 0;
      sum1 = sum1 + arr[z + 1] | 0;
      sum2 = sum2 + arr[z + 2] | 0;
      sum3 = sum3 + arr[z + 3] | 0;
      sum0 = sum0 + arr[z + 4] | 0;
      sum1 = sum1 + arr[z + 5] | 0;
      sum2 = sum2 + arr[z + 6] | 0;
      sum3 = sum3 + arr[z + 7] | 0;
      sum0 = sum0 + arr[z + 8] | 0;
      sum1 = sum1 + arr[z + 9] | 0;
      sum2 = sum2 + arr[z + 10] | 0;
      sum3 = sum3 + arr[z + 11] | 0;
      sum0 = sum0 + arr[z + 12] | 0;
      sum1 = sum1 + arr[z + 13] | 0;
      sum2 = sum2 + arr[z + 14] | 0;
      sum3 = sum3 + arr[z + 15] | 0;
      z += 16;
      N -= 16;
    }
    while (N >= 4) {
      sum0 = sum0 + arr[z + 0] | 0;
      sum1 = sum1 + arr[z + 1] | 0;
      sum2 = sum2 + arr[z + 2] | 0;
      sum3 = sum3 + arr[z + 3] | 0;
      z += 4;
      N -= 4;
    }
    sum3 = ((sum3 + (sum2 << 8) | 0) + (sum1 << 16) | 0) + (sum0 << 24) | 0;
    switch (N) {
      case 3:
        sum3 = sum3 + (arr[z + 2] << 8) | 0;
      case 2:
        sum3 = sum3 + (arr[z + 1] << 16) | 0;
      case 1:
        sum3 = sum3 + (arr[z + 0] << 24) | 0;
    }
    return sum3 >>> 0;
  }
  function applyDelta(source, delta) {
    let total = 0;
    const zDelta = new Reader(delta);
    const lenSrc = source.length;
    const lenDelta = delta.length;
    const limit = zDelta.getInt();
    if (zDelta.getChar() !== "\n")
      throw new Error("size integer not terminated by '\\n'");
    const zOut = new Writer();
    while (zDelta.haveBytes()) {
      const cnt = zDelta.getInt();
      let ofst;
      switch (zDelta.getChar()) {
        case "@":
          ofst = zDelta.getInt();
          if (zDelta.haveBytes() && zDelta.getChar() !== ",")
            throw new Error("copy command not terminated by ','");
          total += cnt;
          if (total > limit)
            throw new Error("copy exceeds output file size");
          if (ofst + cnt > lenSrc)
            throw new Error("copy extends past end of input");
          zOut.putArray(source, ofst, ofst + cnt);
          break;
        case ":":
          total += cnt;
          if (total > limit)
            throw new Error("insert command gives an output larger than predicted");
          if (cnt > lenDelta)
            throw new Error("insert count exceeds size of delta");
          zOut.putArray(zDelta.a, zDelta.pos, zDelta.pos + cnt);
          zDelta.pos += cnt;
          break;
        case ";": {
          const out = zOut.toByteArray(source);
          if (cnt !== checksum(out))
            throw new Error("bad checksum");
          if (total !== limit)
            throw new Error("generated size does not match predicted size");
          return out;
        }
        default:
          throw new Error("unknown delta operator");
      }
    }
    throw new Error("unterminated delta");
  }
  var JsonCodec = class {
    name() {
      return "json";
    }
    encodeCommands(commands) {
      return commands.map((c) => JSON.stringify(c)).join("\n");
    }
    decodeReplies(data2) {
      return data2.trim().split("\n").map((r) => JSON.parse(r));
    }
    applyDeltaIfNeeded(pub, prevValue) {
      let newData, newPrevValue;
      if (pub.delta) {
        const valueArray = applyDelta(prevValue, new TextEncoder().encode(pub.data));
        newData = JSON.parse(new TextDecoder().decode(valueArray));
        newPrevValue = valueArray;
      } else {
        newData = JSON.parse(pub.data);
        newPrevValue = new TextEncoder().encode(pub.data);
      }
      return { newData, newPrevValue };
    }
  };
  var defaults = {
    token: "",
    getToken: null,
    data: null,
    getData: null,
    debug: false,
    name: "js",
    version: "",
    fetch: null,
    readableStream: null,
    websocket: null,
    eventsource: null,
    sockjs: null,
    sockjsOptions: {},
    emulationEndpoint: "/emulation",
    minReconnectDelay: 500,
    maxReconnectDelay: 2e4,
    timeout: 5e3,
    maxServerPingDelay: 1e4,
    networkEventTarget: null
  };
  var UnauthorizedError = class extends Error {
    constructor(message) {
      super(message);
      this.name = this.constructor.name;
    }
  };
  var Centrifuge = class extends EventEmitter$1 {
    /** Constructs Centrifuge client. Call connect() method to start connecting. */
    constructor(endpoint, options) {
      super();
      this._reconnectTimeout = null;
      this._refreshTimeout = null;
      this._serverPingTimeout = null;
      this.state = State.Disconnected;
      this._transportIsOpen = false;
      this._endpoint = endpoint;
      this._emulation = false;
      this._transports = [];
      this._currentTransportIndex = 0;
      this._triedAllTransports = false;
      this._transportWasOpen = false;
      this._transport = null;
      this._transportId = 0;
      this._deviceWentOffline = false;
      this._transportClosed = true;
      this._codec = new JsonCodec();
      this._reconnecting = false;
      this._reconnectTimeout = null;
      this._reconnectAttempts = 0;
      this._client = null;
      this._session = "";
      this._node = "";
      this._subs = {};
      this._serverSubs = {};
      this._commandId = 0;
      this._commands = [];
      this._batching = false;
      this._refreshRequired = false;
      this._refreshTimeout = null;
      this._callbacks = {};
      this._token = "";
      this._data = null;
      this._dispatchPromise = Promise.resolve();
      this._serverPing = 0;
      this._serverPingTimeout = null;
      this._sendPong = false;
      this._promises = {};
      this._promiseId = 0;
      this._debugEnabled = false;
      this._networkEventsSet = false;
      this._config = Object.assign(Object.assign({}, defaults), options);
      this._configure();
      if (this._debugEnabled) {
        this.on("state", (ctx) => {
          this._debug("client state", ctx.oldState, "->", ctx.newState);
        });
        this.on("error", (ctx) => {
          this._debug("client error", ctx);
        });
      } else {
        this.on("error", function() {
          Function.prototype();
        });
      }
    }
    /** newSubscription allocates new Subscription to a channel. Since server only allows
     * one subscription per channel per client this method throws if client already has
     * channel subscription in internal registry.
     * */
    newSubscription(channel, options) {
      if (this.getSubscription(channel) !== null) {
        throw new Error("Subscription to the channel " + channel + " already exists");
      }
      const sub = new Subscription(this, channel, options);
      this._subs[channel] = sub;
      return sub;
    }
    /** getSubscription returns Subscription if it's registered in the internal
     * registry or null. */
    getSubscription(channel) {
      return this._getSub(channel);
    }
    /** removeSubscription allows removing Subcription from the internal registry. Subscrption
     * must be in unsubscribed state. */
    removeSubscription(sub) {
      if (!sub) {
        return;
      }
      if (sub.state !== SubscriptionState.Unsubscribed) {
        sub.unsubscribe();
      }
      this._removeSubscription(sub);
    }
    /** Get a map with all current client-side subscriptions. */
    subscriptions() {
      return this._subs;
    }
    /** ready returns a Promise which resolves upon client goes to Connected
     * state and rejects in case of client goes to Disconnected or Failed state.
     * Users can provide optional timeout in milliseconds. */
    ready(timeout) {
      if (this.state === State.Disconnected) {
        return Promise.reject({ code: errorCodes.clientDisconnected, message: "client disconnected" });
      }
      if (this.state === State.Connected) {
        return Promise.resolve();
      }
      return new Promise((res, rej) => {
        const ctx = {
          resolve: res,
          reject: rej
        };
        if (timeout) {
          ctx.timeout = setTimeout(function() {
            rej({ code: errorCodes.timeout, message: "timeout" });
          }, timeout);
        }
        this._promises[this._nextPromiseId()] = ctx;
      });
    }
    /** connect to a server. */
    connect() {
      if (this._isConnected()) {
        this._debug("connect called when already connected");
        return;
      }
      if (this._isConnecting()) {
        this._debug("connect called when already connecting");
        return;
      }
      this._debug("connect called");
      this._reconnectAttempts = 0;
      this._startConnecting();
    }
    /** disconnect from a server. */
    disconnect() {
      this._disconnect(disconnectedCodes.disconnectCalled, "disconnect called", false);
    }
    /** setToken allows setting connection token. Or resetting used token to be empty.  */
    setToken(token) {
      this._token = token;
    }
    /** send asynchronous data to a server (without any response from a server
     * expected, see rpc method if you need response). */
    send(data2) {
      const cmd = {
        send: {
          data: data2
        }
      };
      const self = this;
      return this._methodCall().then(function() {
        const sent = self._transportSendCommands([cmd]);
        if (!sent) {
          return Promise.reject(self._createErrorObject(errorCodes.transportWriteError, "transport write error"));
        }
        return Promise.resolve();
      });
    }
    /** rpc to a server - i.e. a call which waits for a response with data. */
    rpc(method, data2) {
      const cmd = {
        rpc: {
          method,
          data: data2
        }
      };
      const self = this;
      return this._methodCall().then(function() {
        return self._callPromise(cmd, function(reply) {
          return {
            "data": reply.rpc.data
          };
        });
      });
    }
    /** publish data to a channel. */
    publish(channel, data2) {
      const cmd = {
        publish: {
          channel,
          data: data2
        }
      };
      const self = this;
      return this._methodCall().then(function() {
        return self._callPromise(cmd, function() {
          return {};
        });
      });
    }
    /** history for a channel. By default it does not return publications (only current
     *  StreamPosition data) – provide an explicit limit > 0 to load publications.*/
    history(channel, options) {
      const cmd = {
        history: this._getHistoryRequest(channel, options)
      };
      const self = this;
      return this._methodCall().then(function() {
        return self._callPromise(cmd, function(reply) {
          const result = reply.history;
          const publications = [];
          if (result.publications) {
            for (let i = 0; i < result.publications.length; i++) {
              publications.push(self._getPublicationContext(channel, result.publications[i]));
            }
          }
          return {
            "publications": publications,
            "epoch": result.epoch || "",
            "offset": result.offset || 0
          };
        });
      });
    }
    /** presence for a channel. */
    presence(channel) {
      const cmd = {
        presence: {
          channel
        }
      };
      const self = this;
      return this._methodCall().then(function() {
        return self._callPromise(cmd, function(reply) {
          const clients = reply.presence.presence;
          for (const clientId in clients) {
            if (clients.hasOwnProperty(clientId)) {
              const connInfo = clients[clientId]["conn_info"];
              const chanInfo = clients[clientId]["chan_info"];
              if (connInfo) {
                clients[clientId].connInfo = connInfo;
              }
              if (chanInfo) {
                clients[clientId].chanInfo = chanInfo;
              }
            }
          }
          return {
            "clients": clients
          };
        });
      });
    }
    /** presence stats for a channel. */
    presenceStats(channel) {
      const cmd = {
        "presence_stats": {
          channel
        }
      };
      const self = this;
      return this._methodCall().then(function() {
        return self._callPromise(cmd, function(reply) {
          const result = reply.presence_stats;
          return {
            "numUsers": result.num_users,
            "numClients": result.num_clients
          };
        });
      });
    }
    /** start command batching (collect into temporary buffer without sending to a server)
     * until stopBatching called.*/
    startBatching() {
      this._batching = true;
    }
    /** stop batching commands and flush collected commands to the
     * network (all in one request/frame).*/
    stopBatching() {
      const self = this;
      Promise.resolve().then(function() {
        Promise.resolve().then(function() {
          self._batching = false;
          self._flush();
        });
      });
    }
    _debug(...args) {
      if (!this._debugEnabled) {
        return;
      }
      log("debug", args);
    }
    /** @internal */
    _formatOverride() {
      return;
    }
    _configure() {
      if (!("Promise" in globalThis)) {
        throw new Error("Promise polyfill required");
      }
      if (!this._endpoint) {
        throw new Error("endpoint configuration required");
      }
      if (this._config.token !== null) {
        this._token = this._config.token;
      }
      if (this._config.data !== null) {
        this._data = this._config.data;
      }
      this._codec = new JsonCodec();
      this._formatOverride();
      if (this._config.debug === true || typeof localStorage !== "undefined" && localStorage.getItem("centrifuge.debug")) {
        this._debugEnabled = true;
      }
      this._debug("config", this._config);
      if (typeof this._endpoint === "string") ;
      else if (typeof this._endpoint === "object" && this._endpoint instanceof Array) {
        this._transports = this._endpoint;
        this._emulation = true;
        for (const i in this._transports) {
          if (this._transports.hasOwnProperty(i)) {
            const transportConfig = this._transports[i];
            if (!transportConfig.endpoint || !transportConfig.transport) {
              throw new Error("malformed transport configuration");
            }
            const transportName = transportConfig.transport;
            if (["websocket", "http_stream", "sse", "sockjs", "webtransport"].indexOf(transportName) < 0) {
              throw new Error("unsupported transport name: " + transportName);
            }
          }
        }
      } else {
        throw new Error("unsupported url configuration type: only string or array of objects are supported");
      }
    }
    _setState(newState) {
      if (this.state !== newState) {
        this._reconnecting = false;
        const oldState = this.state;
        this.state = newState;
        this.emit("state", { newState, oldState });
        return true;
      }
      return false;
    }
    _isDisconnected() {
      return this.state === State.Disconnected;
    }
    _isConnecting() {
      return this.state === State.Connecting;
    }
    _isConnected() {
      return this.state === State.Connected;
    }
    _nextCommandId() {
      return ++this._commandId;
    }
    _setNetworkEvents() {
      if (this._networkEventsSet) {
        return;
      }
      let eventTarget = null;
      if (this._config.networkEventTarget !== null) {
        eventTarget = this._config.networkEventTarget;
      } else if (typeof globalThis.addEventListener !== "undefined") {
        eventTarget = globalThis;
      }
      if (eventTarget) {
        eventTarget.addEventListener("offline", () => {
          this._debug("offline event triggered");
          if (this.state === State.Connected || this.state === State.Connecting) {
            this._disconnect(connectingCodes.transportClosed, "transport closed", true);
            this._deviceWentOffline = true;
          }
        });
        eventTarget.addEventListener("online", () => {
          this._debug("online event triggered");
          if (this.state !== State.Connecting) {
            return;
          }
          if (this._deviceWentOffline && !this._transportClosed) {
            this._deviceWentOffline = false;
            this._transportClosed = true;
          }
          this._clearReconnectTimeout();
          this._startReconnecting();
        });
        this._networkEventsSet = true;
      }
    }
    _getReconnectDelay() {
      const delay = backoff(this._reconnectAttempts, this._config.minReconnectDelay, this._config.maxReconnectDelay);
      this._reconnectAttempts += 1;
      return delay;
    }
    _clearOutgoingRequests() {
      for (const id in this._callbacks) {
        if (this._callbacks.hasOwnProperty(id)) {
          const callbacks = this._callbacks[id];
          clearTimeout(callbacks.timeout);
          const errback = callbacks.errback;
          if (!errback) {
            continue;
          }
          errback({ error: this._createErrorObject(errorCodes.connectionClosed, "connection closed") });
        }
      }
      this._callbacks = {};
    }
    _clearConnectedState() {
      this._client = null;
      this._clearServerPingTimeout();
      this._clearRefreshTimeout();
      for (const channel in this._subs) {
        if (!this._subs.hasOwnProperty(channel)) {
          continue;
        }
        const sub = this._subs[channel];
        if (sub.state === SubscriptionState.Subscribed) {
          sub._setSubscribing(subscribingCodes.transportClosed, "transport closed");
        }
      }
      for (const channel in this._serverSubs) {
        if (this._serverSubs.hasOwnProperty(channel)) {
          this.emit("subscribing", { channel });
        }
      }
    }
    _handleWriteError(commands) {
      for (const command of commands) {
        const id = command.id;
        if (!(id in this._callbacks)) {
          continue;
        }
        const callbacks = this._callbacks[id];
        clearTimeout(this._callbacks[id].timeout);
        delete this._callbacks[id];
        const errback = callbacks.errback;
        errback({ error: this._createErrorObject(errorCodes.transportWriteError, "transport write error") });
      }
    }
    _transportSendCommands(commands) {
      if (!commands.length) {
        return true;
      }
      if (!this._transport) {
        return false;
      }
      try {
        this._transport.send(this._codec.encodeCommands(commands), this._session, this._node);
      } catch (e) {
        this._debug("error writing commands", e);
        this._handleWriteError(commands);
        return false;
      }
      return true;
    }
    _initializeTransport() {
      let websocket;
      if (this._config.websocket !== null) {
        websocket = this._config.websocket;
      } else {
        if (!(typeof globalThis.WebSocket !== "function" && typeof globalThis.WebSocket !== "object")) {
          websocket = globalThis.WebSocket;
        }
      }
      let sockjs = null;
      if (this._config.sockjs !== null) {
        sockjs = this._config.sockjs;
      } else {
        if (typeof globalThis.SockJS !== "undefined") {
          sockjs = globalThis.SockJS;
        }
      }
      let eventsource = null;
      if (this._config.eventsource !== null) {
        eventsource = this._config.eventsource;
      } else {
        if (typeof globalThis.EventSource !== "undefined") {
          eventsource = globalThis.EventSource;
        }
      }
      let fetchFunc = null;
      if (this._config.fetch !== null) {
        fetchFunc = this._config.fetch;
      } else {
        if (typeof globalThis.fetch !== "undefined") {
          fetchFunc = globalThis.fetch;
        }
      }
      let readableStream = null;
      if (this._config.readableStream !== null) {
        readableStream = this._config.readableStream;
      } else {
        if (typeof globalThis.ReadableStream !== "undefined") {
          readableStream = globalThis.ReadableStream;
        }
      }
      if (!this._emulation) {
        if (startsWith(this._endpoint, "http")) {
          throw new Error("Provide explicit transport endpoints configuration in case of using HTTP (i.e. using array of TransportEndpoint instead of a single string), or use ws(s):// scheme in an endpoint if you aimed using WebSocket transport");
        } else {
          this._debug("client will use websocket");
          this._transport = new WebsocketTransport(this._endpoint, {
            websocket
          });
          if (!this._transport.supported()) {
            throw new Error("WebSocket not available");
          }
        }
      } else {
        if (this._currentTransportIndex >= this._transports.length) {
          this._triedAllTransports = true;
          this._currentTransportIndex = 0;
        }
        let count = 0;
        while (true) {
          if (count >= this._transports.length) {
            throw new Error("no supported transport found");
          }
          const transportConfig = this._transports[this._currentTransportIndex];
          const transportName = transportConfig.transport;
          const transportEndpoint = transportConfig.endpoint;
          if (transportName === "websocket") {
            this._debug("trying websocket transport");
            this._transport = new WebsocketTransport(transportEndpoint, {
              websocket
            });
            if (!this._transport.supported()) {
              this._debug("websocket transport not available");
              this._currentTransportIndex++;
              count++;
              continue;
            }
          } else if (transportName === "webtransport") {
            this._debug("trying webtransport transport");
            this._transport = new WebtransportTransport(transportEndpoint, {
              webtransport: globalThis.WebTransport,
              decoder: this._codec,
              encoder: this._codec
            });
            if (!this._transport.supported()) {
              this._debug("webtransport transport not available");
              this._currentTransportIndex++;
              count++;
              continue;
            }
          } else if (transportName === "http_stream") {
            this._debug("trying http_stream transport");
            this._transport = new HttpStreamTransport(transportEndpoint, {
              fetch: fetchFunc,
              readableStream,
              emulationEndpoint: this._config.emulationEndpoint,
              decoder: this._codec,
              encoder: this._codec
            });
            if (!this._transport.supported()) {
              this._debug("http_stream transport not available");
              this._currentTransportIndex++;
              count++;
              continue;
            }
          } else if (transportName === "sse") {
            this._debug("trying sse transport");
            this._transport = new SseTransport(transportEndpoint, {
              eventsource,
              fetch: fetchFunc,
              emulationEndpoint: this._config.emulationEndpoint
            });
            if (!this._transport.supported()) {
              this._debug("sse transport not available");
              this._currentTransportIndex++;
              count++;
              continue;
            }
          } else if (transportName === "sockjs") {
            this._debug("trying sockjs");
            this._transport = new SockjsTransport(transportEndpoint, {
              sockjs,
              sockjsOptions: this._config.sockjsOptions
            });
            if (!this._transport.supported()) {
              this._debug("sockjs transport not available");
              this._currentTransportIndex++;
              count++;
              continue;
            }
          } else {
            throw new Error("unknown transport " + transportName);
          }
          break;
        }
      }
      const self = this;
      const transport = this._transport;
      const transportId = this._nextTransportId();
      self._debug("id of transport", transportId);
      let wasOpen = false;
      const initialCommands = [];
      if (this._transport.emulation()) {
        const connectCommand = self._sendConnect(true);
        initialCommands.push(connectCommand);
      }
      this._setNetworkEvents();
      const initialData = this._codec.encodeCommands(initialCommands);
      this._transportClosed = false;
      let connectTimeout;
      connectTimeout = setTimeout(function() {
        transport.close();
      }, this._config.timeout);
      this._transport.initialize(this._codec.name(), {
        onOpen: function() {
          if (connectTimeout) {
            clearTimeout(connectTimeout);
            connectTimeout = null;
          }
          if (self._transportId != transportId) {
            self._debug("open callback from non-actual transport");
            transport.close();
            return;
          }
          wasOpen = true;
          self._debug(transport.subName(), "transport open");
          if (transport.emulation()) {
            return;
          }
          self._transportIsOpen = true;
          self._transportWasOpen = true;
          self.startBatching();
          self._sendConnect(false);
          self._sendSubscribeCommands();
          self.stopBatching();
          self.emit("__centrifuge_debug:connect_frame_sent", {});
        },
        onError: function(e) {
          if (self._transportId != transportId) {
            self._debug("error callback from non-actual transport");
            return;
          }
          self._debug("transport level error", e);
        },
        onClose: function(closeEvent) {
          if (connectTimeout) {
            clearTimeout(connectTimeout);
            connectTimeout = null;
          }
          if (self._transportId != transportId) {
            self._debug("close callback from non-actual transport");
            return;
          }
          self._debug(transport.subName(), "transport closed");
          self._transportClosed = true;
          self._transportIsOpen = false;
          let reason = "connection closed";
          let needReconnect = true;
          let code = 0;
          if (closeEvent && "code" in closeEvent && closeEvent.code) {
            code = closeEvent.code;
          }
          if (closeEvent && closeEvent.reason) {
            try {
              const advice = JSON.parse(closeEvent.reason);
              reason = advice.reason;
              needReconnect = advice.reconnect;
            } catch (e) {
              reason = closeEvent.reason;
              if (code >= 3500 && code < 4e3 || code >= 4500 && code < 5e3) {
                needReconnect = false;
              }
            }
          }
          if (code < 3e3) {
            if (code === 1009) {
              code = disconnectedCodes.messageSizeLimit;
              reason = "message size limit exceeded";
              needReconnect = false;
            } else {
              code = connectingCodes.transportClosed;
              reason = "transport closed";
            }
            if (self._emulation && !self._transportWasOpen) {
              self._currentTransportIndex++;
              if (self._currentTransportIndex >= self._transports.length) {
                self._triedAllTransports = true;
                self._currentTransportIndex = 0;
              }
            }
          } else {
            self._transportWasOpen = true;
          }
          if (self._isConnecting() && !wasOpen) {
            self.emit("error", {
              type: "transport",
              error: {
                code: errorCodes.transportClosed,
                message: "transport closed"
              },
              transport: transport.name()
            });
          }
          self._reconnecting = false;
          self._disconnect(code, reason, needReconnect);
        },
        onMessage: function(data2) {
          self._dataReceived(data2);
        }
      }, initialData);
      self.emit("__centrifuge_debug:transport_initialized", {});
    }
    _sendConnect(skipSending) {
      const connectCommand = this._constructConnectCommand();
      const self = this;
      this._call(connectCommand, skipSending).then((resolveCtx) => {
        const result = resolveCtx.reply.connect;
        self._connectResponse(result);
        if (resolveCtx.next) {
          resolveCtx.next();
        }
      }, (rejectCtx) => {
        self._connectError(rejectCtx.error);
        if (rejectCtx.next) {
          rejectCtx.next();
        }
      });
      return connectCommand;
    }
    _startReconnecting() {
      this._debug("start reconnecting");
      if (!this._isConnecting()) {
        this._debug("stop reconnecting: client not in connecting state");
        return;
      }
      if (this._reconnecting) {
        this._debug("reconnect already in progress, return from reconnect routine");
        return;
      }
      if (this._transportClosed === false) {
        this._debug("waiting for transport close");
        return;
      }
      this._reconnecting = true;
      const self = this;
      const emptyToken = this._token === "";
      const needTokenRefresh = this._refreshRequired || emptyToken && this._config.getToken !== null;
      if (!needTokenRefresh) {
        if (this._config.getData) {
          this._config.getData().then(function(data2) {
            if (!self._isConnecting()) {
              return;
            }
            self._data = data2;
            self._initializeTransport();
          });
        } else {
          this._initializeTransport();
        }
        return;
      }
      this._getToken().then(function(token) {
        if (!self._isConnecting()) {
          return;
        }
        if (token == null || token == void 0) {
          self._failUnauthorized();
          return;
        }
        self._token = token;
        self._debug("connection token refreshed");
        if (self._config.getData) {
          self._config.getData().then(function(data2) {
            if (!self._isConnecting()) {
              return;
            }
            self._data = data2;
            self._initializeTransport();
          });
        } else {
          self._initializeTransport();
        }
      }).catch(function(e) {
        if (!self._isConnecting()) {
          return;
        }
        if (e instanceof UnauthorizedError) {
          self._failUnauthorized();
          return;
        }
        self.emit("error", {
          "type": "connectToken",
          "error": {
            code: errorCodes.clientConnectToken,
            message: e !== void 0 ? e.toString() : ""
          }
        });
        const delay = self._getReconnectDelay();
        self._debug("error on connection token refresh, reconnect after " + delay + " milliseconds", e);
        self._reconnecting = false;
        self._reconnectTimeout = setTimeout(() => {
          self._startReconnecting();
        }, delay);
      });
    }
    _connectError(err) {
      if (this.state !== State.Connecting) {
        return;
      }
      if (err.code === 109) {
        this._refreshRequired = true;
      }
      if (err.code < 100 || err.temporary === true || err.code === 109) {
        this.emit("error", {
          "type": "connect",
          "error": err
        });
        this._debug("closing transport due to connect error");
        this._disconnect(err.code, err.message, true);
      } else {
        this._disconnect(err.code, err.message, false);
      }
    }
    _scheduleReconnect() {
      if (!this._isConnecting()) {
        return;
      }
      let isInitialHandshake = false;
      if (this._emulation && !this._transportWasOpen && !this._triedAllTransports) {
        isInitialHandshake = true;
      }
      let delay = this._getReconnectDelay();
      if (isInitialHandshake) {
        delay = 0;
      }
      this._debug("reconnect after " + delay + " milliseconds");
      this._clearReconnectTimeout();
      this._reconnectTimeout = setTimeout(() => {
        this._startReconnecting();
      }, delay);
    }
    _constructConnectCommand() {
      const req = {};
      if (this._token) {
        req.token = this._token;
      }
      if (this._data) {
        req.data = this._data;
      }
      if (this._config.name) {
        req.name = this._config.name;
      }
      if (this._config.version) {
        req.version = this._config.version;
      }
      const subs = {};
      let hasSubs = false;
      for (const channel in this._serverSubs) {
        if (this._serverSubs.hasOwnProperty(channel) && this._serverSubs[channel].recoverable) {
          hasSubs = true;
          const sub = {
            "recover": true
          };
          if (this._serverSubs[channel].offset) {
            sub["offset"] = this._serverSubs[channel].offset;
          }
          if (this._serverSubs[channel].epoch) {
            sub["epoch"] = this._serverSubs[channel].epoch;
          }
          subs[channel] = sub;
        }
      }
      if (hasSubs) {
        req.subs = subs;
      }
      return {
        connect: req
      };
    }
    _getHistoryRequest(channel, options) {
      const req = {
        channel
      };
      if (options !== void 0) {
        if (options.since) {
          req.since = {
            offset: options.since.offset
          };
          if (options.since.epoch) {
            req.since.epoch = options.since.epoch;
          }
        }
        if (options.limit !== void 0) {
          req.limit = options.limit;
        }
        if (options.reverse === true) {
          req.reverse = true;
        }
      }
      return req;
    }
    _methodCall() {
      if (this._isConnected()) {
        return Promise.resolve();
      }
      return new Promise((res, rej) => {
        const timeout = setTimeout(function() {
          rej({ code: errorCodes.timeout, message: "timeout" });
        }, this._config.timeout);
        this._promises[this._nextPromiseId()] = {
          timeout,
          resolve: res,
          reject: rej
        };
      });
    }
    _callPromise(cmd, resultCB) {
      return new Promise((resolve, reject) => {
        this._call(cmd, false).then((resolveCtx) => {
          resolve(resultCB(resolveCtx.reply));
          if (resolveCtx.next) {
            resolveCtx.next();
          }
        }, (rejectCtx) => {
          reject(rejectCtx.error);
          if (rejectCtx.next) {
            rejectCtx.next();
          }
        });
      });
    }
    _dataReceived(data2) {
      if (this._serverPing > 0) {
        this._waitServerPing();
      }
      const replies = this._codec.decodeReplies(data2);
      this._dispatchPromise = this._dispatchPromise.then(() => {
        let finishDispatch;
        this._dispatchPromise = new Promise((resolve) => {
          finishDispatch = resolve;
        });
        this._dispatchSynchronized(replies, finishDispatch);
      });
    }
    _dispatchSynchronized(replies, finishDispatch) {
      let p = Promise.resolve();
      for (const i in replies) {
        if (replies.hasOwnProperty(i)) {
          p = p.then(() => {
            return this._dispatchReply(replies[i]);
          });
        }
      }
      p = p.then(() => {
        finishDispatch();
      });
    }
    _dispatchReply(reply) {
      let next;
      const p = new Promise((resolve) => {
        next = resolve;
      });
      if (reply === void 0 || reply === null) {
        this._debug("dispatch: got undefined or null reply");
        next();
        return p;
      }
      const id = reply.id;
      if (id && id > 0) {
        this._handleReply(reply, next);
      } else {
        if (!reply.push) {
          this._handleServerPing(next);
        } else {
          this._handlePush(reply.push, next);
        }
      }
      return p;
    }
    _call(cmd, skipSending) {
      return new Promise((resolve, reject) => {
        cmd.id = this._nextCommandId();
        this._registerCall(cmd.id, resolve, reject);
        if (!skipSending) {
          this._addCommand(cmd);
        }
      });
    }
    _startConnecting() {
      this._debug("start connecting");
      if (this._setState(State.Connecting)) {
        this.emit("connecting", { code: connectingCodes.connectCalled, reason: "connect called" });
      }
      this._client = null;
      this._startReconnecting();
    }
    _disconnect(code, reason, reconnect) {
      if (this._isDisconnected()) {
        return;
      }
      this._transportIsOpen = false;
      const previousState = this.state;
      this._reconnecting = false;
      const ctx = {
        code,
        reason
      };
      let needEvent = false;
      if (reconnect) {
        needEvent = this._setState(State.Connecting);
      } else {
        needEvent = this._setState(State.Disconnected);
        this._rejectPromises({ code: errorCodes.clientDisconnected, message: "disconnected" });
      }
      this._clearOutgoingRequests();
      if (previousState === State.Connecting) {
        this._clearReconnectTimeout();
      }
      if (previousState === State.Connected) {
        this._clearConnectedState();
      }
      if (needEvent) {
        if (this._isConnecting()) {
          this.emit("connecting", ctx);
        } else {
          this.emit("disconnected", ctx);
        }
      }
      if (this._transport) {
        this._debug("closing existing transport");
        const transport = this._transport;
        this._transport = null;
        transport.close();
        this._transportClosed = true;
        this._nextTransportId();
      } else {
        this._debug("no transport to close");
      }
      this._scheduleReconnect();
    }
    _failUnauthorized() {
      this._disconnect(disconnectedCodes.unauthorized, "unauthorized", false);
    }
    _getToken() {
      this._debug("get connection token");
      if (!this._config.getToken) {
        this.emit("error", {
          type: "configuration",
          error: {
            code: errorCodes.badConfiguration,
            message: "token expired but no getToken function set in the configuration"
          }
        });
        throw new UnauthorizedError("");
      }
      return this._config.getToken({});
    }
    _refresh() {
      const clientId = this._client;
      const self = this;
      this._getToken().then(function(token) {
        if (clientId !== self._client) {
          return;
        }
        if (!token) {
          self._failUnauthorized();
          return;
        }
        self._token = token;
        self._debug("connection token refreshed");
        if (!self._isConnected()) {
          return;
        }
        const cmd = {
          refresh: { token: self._token }
        };
        self._call(cmd, false).then((resolveCtx) => {
          const result = resolveCtx.reply.refresh;
          self._refreshResponse(result);
          if (resolveCtx.next) {
            resolveCtx.next();
          }
        }, (rejectCtx) => {
          self._refreshError(rejectCtx.error);
          if (rejectCtx.next) {
            rejectCtx.next();
          }
        });
      }).catch(function(e) {
        if (!self._isConnected()) {
          return;
        }
        if (e instanceof UnauthorizedError) {
          self._failUnauthorized();
          return;
        }
        self.emit("error", {
          type: "refreshToken",
          error: {
            code: errorCodes.clientRefreshToken,
            message: e !== void 0 ? e.toString() : ""
          }
        });
        self._refreshTimeout = setTimeout(() => self._refresh(), self._getRefreshRetryDelay());
      });
    }
    _refreshError(err) {
      if (err.code < 100 || err.temporary === true) {
        this.emit("error", {
          type: "refresh",
          error: err
        });
        this._refreshTimeout = setTimeout(() => this._refresh(), this._getRefreshRetryDelay());
      } else {
        this._disconnect(err.code, err.message, false);
      }
    }
    _getRefreshRetryDelay() {
      return backoff(0, 5e3, 1e4);
    }
    _refreshResponse(result) {
      if (this._refreshTimeout) {
        clearTimeout(this._refreshTimeout);
        this._refreshTimeout = null;
      }
      if (result.expires) {
        this._client = result.client;
        this._refreshTimeout = setTimeout(() => this._refresh(), ttlMilliseconds(result.ttl));
      }
    }
    _removeSubscription(sub) {
      if (sub === null) {
        return;
      }
      delete this._subs[sub.channel];
    }
    _unsubscribe(sub) {
      if (!this._transportIsOpen) {
        return Promise.resolve();
      }
      const req = {
        channel: sub.channel
      };
      const cmd = { unsubscribe: req };
      const self = this;
      const unsubscribePromise = new Promise((resolve, _) => {
        this._call(cmd, false).then((resolveCtx) => {
          resolve();
          if (resolveCtx.next) {
            resolveCtx.next();
          }
        }, (rejectCtx) => {
          resolve();
          if (rejectCtx.next) {
            rejectCtx.next();
          }
          self._disconnect(connectingCodes.unsubscribeError, "unsubscribe error", true);
        });
      });
      return unsubscribePromise;
    }
    _getSub(channel) {
      const sub = this._subs[channel];
      if (!sub) {
        return null;
      }
      return sub;
    }
    _isServerSub(channel) {
      return this._serverSubs[channel] !== void 0;
    }
    _sendSubscribeCommands() {
      const commands = [];
      for (const channel in this._subs) {
        if (!this._subs.hasOwnProperty(channel)) {
          continue;
        }
        const sub = this._subs[channel];
        if (sub._inflight === true) {
          continue;
        }
        if (sub.state === SubscriptionState.Subscribing) {
          const cmd = sub._subscribe();
          if (cmd) {
            commands.push(cmd);
          }
        }
      }
      return commands;
    }
    _connectResponse(result) {
      this._transportIsOpen = true;
      this._transportWasOpen = true;
      this._reconnectAttempts = 0;
      this._refreshRequired = false;
      if (this._isConnected()) {
        return;
      }
      this._client = result.client;
      this._setState(State.Connected);
      if (this._refreshTimeout) {
        clearTimeout(this._refreshTimeout);
      }
      if (result.expires) {
        this._refreshTimeout = setTimeout(() => this._refresh(), ttlMilliseconds(result.ttl));
      }
      this._session = result.session;
      this._node = result.node;
      this.startBatching();
      this._sendSubscribeCommands();
      this.stopBatching();
      const ctx = {
        client: result.client,
        transport: this._transport.subName()
      };
      if (result.data) {
        ctx.data = result.data;
      }
      this.emit("connected", ctx);
      this._resolvePromises();
      this._processServerSubs(result.subs || {});
      if (result.ping && result.ping > 0) {
        this._serverPing = result.ping * 1e3;
        this._sendPong = result.pong === true;
        this._waitServerPing();
      } else {
        this._serverPing = 0;
      }
    }
    _processServerSubs(subs) {
      for (const channel in subs) {
        if (!subs.hasOwnProperty(channel)) {
          continue;
        }
        const sub = subs[channel];
        this._serverSubs[channel] = {
          "offset": sub.offset,
          "epoch": sub.epoch,
          "recoverable": sub.recoverable || false
        };
        const subCtx = this._getSubscribeContext(channel, sub);
        this.emit("subscribed", subCtx);
      }
      for (const channel in subs) {
        if (!subs.hasOwnProperty(channel)) {
          continue;
        }
        const sub = subs[channel];
        if (sub.recovered) {
          const pubs = sub.publications;
          if (pubs && pubs.length > 0) {
            for (const i in pubs) {
              if (pubs.hasOwnProperty(i)) {
                this._handlePublication(channel, pubs[i]);
              }
            }
          }
        }
      }
      for (const channel in this._serverSubs) {
        if (!this._serverSubs.hasOwnProperty(channel)) {
          continue;
        }
        if (!subs[channel]) {
          this.emit("unsubscribed", { channel });
          delete this._serverSubs[channel];
        }
      }
    }
    _clearRefreshTimeout() {
      if (this._refreshTimeout !== null) {
        clearTimeout(this._refreshTimeout);
        this._refreshTimeout = null;
      }
    }
    _clearReconnectTimeout() {
      if (this._reconnectTimeout !== null) {
        clearTimeout(this._reconnectTimeout);
        this._reconnectTimeout = null;
      }
    }
    _clearServerPingTimeout() {
      if (this._serverPingTimeout !== null) {
        clearTimeout(this._serverPingTimeout);
        this._serverPingTimeout = null;
      }
    }
    _waitServerPing() {
      if (this._config.maxServerPingDelay === 0) {
        return;
      }
      if (!this._isConnected()) {
        return;
      }
      this._clearServerPingTimeout();
      this._serverPingTimeout = setTimeout(() => {
        if (!this._isConnected()) {
          return;
        }
        this._disconnect(connectingCodes.noPing, "no ping", true);
      }, this._serverPing + this._config.maxServerPingDelay);
    }
    _getSubscribeContext(channel, result) {
      const ctx = {
        channel,
        positioned: false,
        recoverable: false,
        wasRecovering: false,
        recovered: false
      };
      if (result.recovered) {
        ctx.recovered = true;
      }
      if (result.positioned) {
        ctx.positioned = true;
      }
      if (result.recoverable) {
        ctx.recoverable = true;
      }
      if (result.was_recovering) {
        ctx.wasRecovering = true;
      }
      let epoch = "";
      if ("epoch" in result) {
        epoch = result.epoch;
      }
      let offset = 0;
      if ("offset" in result) {
        offset = result.offset;
      }
      if (ctx.positioned || ctx.recoverable) {
        ctx.streamPosition = {
          "offset": offset,
          "epoch": epoch
        };
      }
      if (result.data) {
        ctx.data = result.data;
      }
      return ctx;
    }
    _handleReply(reply, next) {
      const id = reply.id;
      if (!(id in this._callbacks)) {
        next();
        return;
      }
      const callbacks = this._callbacks[id];
      clearTimeout(this._callbacks[id].timeout);
      delete this._callbacks[id];
      if (!errorExists(reply)) {
        const callback = callbacks.callback;
        if (!callback) {
          return;
        }
        callback({ reply, next });
      } else {
        const errback = callbacks.errback;
        if (!errback) {
          next();
          return;
        }
        const error2 = reply.error;
        errback({ error: error2, next });
      }
    }
    _handleJoin(channel, join) {
      const sub = this._getSub(channel);
      if (!sub) {
        if (this._isServerSub(channel)) {
          const ctx = { channel, info: this._getJoinLeaveContext(join.info) };
          this.emit("join", ctx);
        }
        return;
      }
      sub._handleJoin(join);
    }
    _handleLeave(channel, leave) {
      const sub = this._getSub(channel);
      if (!sub) {
        if (this._isServerSub(channel)) {
          const ctx = { channel, info: this._getJoinLeaveContext(leave.info) };
          this.emit("leave", ctx);
        }
        return;
      }
      sub._handleLeave(leave);
    }
    _handleUnsubscribe(channel, unsubscribe) {
      const sub = this._getSub(channel);
      if (!sub) {
        if (this._isServerSub(channel)) {
          delete this._serverSubs[channel];
          this.emit("unsubscribed", { channel });
        }
        return;
      }
      if (unsubscribe.code < 2500) {
        sub._setUnsubscribed(unsubscribe.code, unsubscribe.reason, false);
      } else {
        sub._setSubscribing(unsubscribe.code, unsubscribe.reason);
      }
    }
    _handleSubscribe(channel, sub) {
      this._serverSubs[channel] = {
        "offset": sub.offset,
        "epoch": sub.epoch,
        "recoverable": sub.recoverable || false
      };
      this.emit("subscribed", this._getSubscribeContext(channel, sub));
    }
    _handleDisconnect(disconnect) {
      const code = disconnect.code;
      let reconnect = true;
      if (code >= 3500 && code < 4e3 || code >= 4500 && code < 5e3) {
        reconnect = false;
      }
      this._disconnect(code, disconnect.reason, reconnect);
    }
    _getPublicationContext(channel, pub) {
      const ctx = {
        channel,
        data: pub.data
      };
      if (pub.offset) {
        ctx.offset = pub.offset;
      }
      if (pub.info) {
        ctx.info = this._getJoinLeaveContext(pub.info);
      }
      if (pub.tags) {
        ctx.tags = pub.tags;
      }
      return ctx;
    }
    _getJoinLeaveContext(clientInfo) {
      const info = {
        client: clientInfo.client,
        user: clientInfo.user
      };
      if (clientInfo.conn_info) {
        info.connInfo = clientInfo.conn_info;
      }
      if (clientInfo.chan_info) {
        info.chanInfo = clientInfo.chan_info;
      }
      return info;
    }
    _handlePublication(channel, pub) {
      const sub = this._getSub(channel);
      if (!sub) {
        if (this._isServerSub(channel)) {
          const ctx = this._getPublicationContext(channel, pub);
          this.emit("publication", ctx);
          if (pub.offset !== void 0) {
            this._serverSubs[channel].offset = pub.offset;
          }
        }
        return;
      }
      sub._handlePublication(pub);
    }
    _handleMessage(message) {
      this.emit("message", { data: message.data });
    }
    _handleServerPing(next) {
      if (this._sendPong) {
        const cmd = {};
        this._transportSendCommands([cmd]);
      }
      next();
    }
    _handlePush(data2, next) {
      const channel = data2.channel;
      if (data2.pub) {
        this._handlePublication(channel, data2.pub);
      } else if (data2.message) {
        this._handleMessage(data2.message);
      } else if (data2.join) {
        this._handleJoin(channel, data2.join);
      } else if (data2.leave) {
        this._handleLeave(channel, data2.leave);
      } else if (data2.unsubscribe) {
        this._handleUnsubscribe(channel, data2.unsubscribe);
      } else if (data2.subscribe) {
        this._handleSubscribe(channel, data2.subscribe);
      } else if (data2.disconnect) {
        this._handleDisconnect(data2.disconnect);
      }
      next();
    }
    _flush() {
      const commands = this._commands.slice(0);
      this._commands = [];
      this._transportSendCommands(commands);
    }
    _createErrorObject(code, message, temporary) {
      const errObject = {
        code,
        message
      };
      if (temporary) {
        errObject.temporary = true;
      }
      return errObject;
    }
    _registerCall(id, callback, errback) {
      this._callbacks[id] = {
        callback,
        errback,
        timeout: null
      };
      this._callbacks[id].timeout = setTimeout(() => {
        delete this._callbacks[id];
        if (isFunction(errback)) {
          errback({ error: this._createErrorObject(errorCodes.timeout, "timeout") });
        }
      }, this._config.timeout);
    }
    _addCommand(command) {
      if (this._batching) {
        this._commands.push(command);
      } else {
        this._transportSendCommands([command]);
      }
    }
    _nextPromiseId() {
      return ++this._promiseId;
    }
    _nextTransportId() {
      return ++this._transportId;
    }
    _resolvePromises() {
      for (const id in this._promises) {
        if (!this._promises.hasOwnProperty(id)) {
          continue;
        }
        if (this._promises[id].timeout) {
          clearTimeout(this._promises[id].timeout);
        }
        this._promises[id].resolve();
        delete this._promises[id];
      }
    }
    _rejectPromises(err) {
      for (const id in this._promises) {
        if (!this._promises.hasOwnProperty(id)) {
          continue;
        }
        if (this._promises[id].timeout) {
          clearTimeout(this._promises[id].timeout);
        }
        this._promises[id].reject(err);
        delete this._promises[id];
      }
    }
  };
  Centrifuge.SubscriptionState = SubscriptionState;
  Centrifuge.State = State;
  Centrifuge.UnauthorizedError = UnauthorizedError;

  // bff/client/index.ts
  (function() {
    window.Alpine = module_default;
    module_default.store("questionData", {
      async init() {
        this.user = await (await fetch("/user")).json();
        initCentrifugo(this.user);
      },
      questions: JSON.parse(document.getElementById("questions").textContent),
      get sortedQuestions() {
        return this.questions.sort((a, b) => b.Votes - a.Votes);
      },
      user: null,
      addQuestion(question) {
        this.questions.push(question);
      }
    });
    module_default.start();
    document.addEventListener("alpine:init", () => {
    });
    const initCentrifugo = async (user) => {
      console.log("init centrifugo");
      const centrifuge = new Centrifuge("ws://localhost:3333/api/v1/connection/websocket", {
        token: user.Token
      });
      centrifuge.on("connecting", function(ctx) {
        console.log(`connecting: ${ctx.code}, ${ctx.reason}`);
      }).on("connected", function(ctx) {
        console.log(`connected over ${ctx.transport}`);
      }).on("disconnected", function(ctx) {
        console.log(`disconnected: ${ctx.code}, ${ctx.reason}`);
      }).on("message", function(msg) {
        console.log(`message: ${JSON.stringify(msg)}`);
        const data2 = JSON.parse(msg.data.Payload);
        const eventType = msg.data.EventType;
        switch (eventType) {
          case "start_session":
            break;
          case "stop_session":
            break;
          case "user_connected":
          case "user_disconnected":
            break;
          case "new_question":
            module_default.store("questionData").addQuestion(data2);
          case "undo_upvote_question":
            break;
          case "update_question":
            break;
          case "delete_question":
            break;
          case "answer_question":
            break;
        }
      }).connect();
    };
  })();
})();
