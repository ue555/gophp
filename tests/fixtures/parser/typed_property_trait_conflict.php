<?php
echo "before";
trait TypedTraitConflict { public int $value = 1; }
class TypedTraitConflictUser {
    use TypedTraitConflict;
    public string $value = 'text';
}
echo "after";
