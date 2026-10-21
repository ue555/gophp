<?php
function truthyUnion(true|int|null $value): int|true|null { return $value; }
echo gettype(truthyUnion(true)), '|', gettype(truthyUnion(7)), '|', gettype(truthyUnion(null)), "\n";
