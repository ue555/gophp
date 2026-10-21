<?php
namespace DnfExample;
interface A {}
interface B {}
class Both implements A, B {}
use DnfExample\A as First;
use DnfExample\B as Second;
function take((First&Second)|null $x): (Second&First)|null { return $x; }
echo gettype(take(new Both)), '|', gettype(take(null)), "\n";
