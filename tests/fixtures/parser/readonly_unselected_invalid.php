<?php
echo 'before';
if (false) {
    class UnselectedInvalidReadonly { public readonly $value; }
}
