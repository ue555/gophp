<?php
echo "before";
class TypedInvariantParent { public int $value; }
class TypedInvariantChild extends TypedInvariantParent { public string $value; }
echo "after";
