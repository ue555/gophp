<?php
function makeError($value) { return new Exception('stored'); }
$error = makeError([1, 2]);
echo $error->getLine(), '|', count($error->getTrace()), '|';
echo $error->getTrace()[0]['function'], '|';
echo gettype($error->getTrace()[0]['args'][0]), '|';
try { throw $error; } catch (Exception $caught) {
    echo ($caught === $error ? 'same' : 'different'), '|';
    echo $caught->getLine(), '|', count($caught->getTrace()), "\n";
}
