<?php
try {
    new CachingIterator(new stdClass);
    echo 'bad';
} catch (TypeError $e) {
    echo 'type|';
}
$cache = new CachingIterator(new ArrayIterator([1]));
try {
    $cache->getCache();
    echo 'bad';
} catch (BadMethodCallException $e) {
    echo 'cache|';
}
try {
    $cache->setFlags(CachingIterator::CALL_TOSTRING | CachingIterator::TOSTRING_USE_KEY);
    echo 'bad';
} catch (ValueError $e) {
    echo 'flags';
}
