<?php
function maybeUnion(int|false|null $value): int|false|null { return $value; }
echo gettype(maybeUnion(null)), '|', gettype(maybeUnion(false)), '|', gettype(maybeUnion(8)), "\n";
