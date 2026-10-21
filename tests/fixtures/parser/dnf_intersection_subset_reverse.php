<?php
interface A {}
interface B {}
interface C {}
echo "before|";
function redundantReverse((A&B&C)|(A&B) $value) {}
echo "after";
