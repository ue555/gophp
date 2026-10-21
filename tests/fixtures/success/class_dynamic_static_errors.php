<?php
class DynamicError { public function instance() { return 1; } }
$class = 12;
try { $class::missing(); } catch (Error $e) { echo 'class|'; }
$class = 'DynamicError';
try { $class::missing(); } catch (Error $e) { echo 'method|'; }
try { $class::instance(); } catch (Error $e) { echo 'static|'; }
$method = 42;
try { $class::$method(); } catch (Error $e) { echo 'method-type'; }
echo "\n";
