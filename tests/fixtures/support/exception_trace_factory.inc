<?php
function exceptionTraceFactory($input) {
    return new Exception('included');
}
return exceptionTraceFactory(3);
