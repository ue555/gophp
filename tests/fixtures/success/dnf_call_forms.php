<?php
interface CallA {}
interface CallB {}
class CallBoth implements CallA, CallB {}
class CallBox { public static function take((CallA&CallB)|int $x): (CallB&CallA)|int { return $x; } }
$closure = function ((CallA&CallB)|int $x): (CallB&CallA)|int { return $x; };
$arrow = fn ((CallA&CallB)|int $x): (CallB&CallA)|int => $x;
function callVariadic((CallA&CallB)|int ...$xs): int { return count($xs); }
echo gettype(CallBox::take(x: new CallBoth)), '|', gettype($closure(4)), '|', gettype($arrow(new CallBoth)), '|', callVariadic(...[new CallBoth, 3]), "\n";
