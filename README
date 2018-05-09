rfcat library for Go
====================

This package is functional but not very polished. A patch needs to be made to github.com/google/gousb in order to avoid a lockup at startup:

    diff --git a/config.go b/config.go
    index 70e2330..5eb0f76 100644
    --- a/config.go
    +++ b/config.go
    @@ -124,9 +124,11 @@ func (c *Config) Interface(num, alt int) (*Interface, error) {
            return nil, fmt.Errorf("failed to claim interface %d on %s: %v", num, c, err)
        }

    -	if err := c.dev.ctx.libusb.setAlt(c.dev.handle, uint8(num), uint8(alt)); err != nil {
    -		c.dev.ctx.libusb.release(c.dev.handle, uint8(num))
    -		return nil, fmt.Errorf("failed to set alternate config %d on interface %d of %s: %v", alt, num, c, err)
    +	if alt != 0 {
    +		if err := c.dev.ctx.libusb.setAlt(c.dev.handle, uint8(num), uint8(alt)); err != nil {
    +			c.dev.ctx.libusb.release(c.dev.handle, uint8(num))
    +			return nil, fmt.Errorf("failed to set alternate config %d on interface %d of %s: %v", alt, num, c, err)
    +		}
        }

        c.claimed[num] = true
