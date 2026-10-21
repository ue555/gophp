<?php
echo "before";
function &invalidReferenceGenerator() {
    yield from [1, 2];
}
echo "after";
