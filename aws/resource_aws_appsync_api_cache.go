package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appsync"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAwsAppsyncApiCache() *schema.Resource {

	return &schema.Resource{
		Create: resourceAwsAppsyncApiCacheCreate,
		Read:   resourceAwsAppsyncApiCacheRead,
		Update: resourceAwsAppsyncApiCacheUpdate,
		Delete: resourceAwsAppsyncApiCacheDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"api_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_caching_behavior": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					appsync.ApiCachingBehaviorFullRequestCaching,
					appsync.ApiCachingBehaviorPerResolverCaching,
				}, true),
			},
			"at_rest_encryption_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"transit_encryption_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 3600),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					appsync.ApiCacheTypeSmall,
					appsync.ApiCacheTypeMedium,
					appsync.ApiCacheTypeLarge,
					appsync.ApiCacheTypeXlarge,
					appsync.ApiCacheTypeLarge2x,
					appsync.ApiCacheTypeLarge4x,
					appsync.ApiCacheTypeLarge8x,
					appsync.ApiCacheTypeLarge12x,
				}, true),
			},
		},
	}
}

func resourceAwsAppsyncApiCacheCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appsyncconn

	apiID := d.Get("api_id").(string)

	params := &appsync.CreateApiCacheInput{
		ApiCachingBehavior: aws.String(d.Get("api_caching_behavior").(string)),
		ApiId:              aws.String(apiID),
		Ttl:                aws.Int64(d.Get("ttl").(int64)),
	}
	if v, ok := d.GetOk("at_rest_encryption_enabled"); ok {
		params.AtRestEncryptionEnabled = aws.Bool(v.(bool))
	}
	if v, ok := d.GetOk("transit_encryption_enabled"); ok {
		params.TransitEncryptionEnabled = aws.Bool(v.(bool))
	}
	_, err := conn.CreateApiCache(params)
	if err != nil {
		return fmt.Errorf("error creating Appsync API Cache: %s", err)
	}

	d.SetId(fmt.Sprintf("%s", apiID))
	return resourceAwsAppsyncApiCacheRead(d, meta)
}

func resourceAwsAppsyncApiCacheRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appsyncconn

	apiID := d.Id()

	params := &appsync.GetApiCacheInput{
		ApiId: aws.String(apiID),
	}

	resp, err := conn.GetApiCache(params)
	if err != nil {
		return fmt.Errorf("error getting Appsync API Cache for API ID %q: %s", d.Id(), err)
	}
	if resp == nil {
		log.Printf("[WARN] AppSync API Cache for API ID %q not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("api_id", apiID)
	d.Set("api_caching_behavior", resp.ApiCache.ApiCachingBehavior)
	d.Set("ttl", resp.ApiCache.Ttl)
	d.Set("type", resp.ApiCache.Type)

	if err := d.Set("at_rest_encryption_enabled", aws.BoolValue(resp.ApiCache.AtRestEncryptionEnabled)); err != nil {
		return fmt.Errorf("error setting at_rest_encryption_enabled: %s", err)
	}
	if err := d.Set("transit_encryption_enabled", aws.BoolValue(resp.ApiCache.TransitEncryptionEnabled)); err != nil {
		return fmt.Errorf("error setting transit_encryption_enabled: %s", err)
	}

	return nil
}

func resourceAwsAppsyncApiCacheUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appsyncconn

	apiID := d.Id()

	params := &appsync.UpdateApiCacheInput{
		ApiId: aws.String(apiID),
	}
	if d.HasChange("api_caching_behavior") {
		params.ApiCachingBehavior = aws.String(d.Get("api_caching_behavior").(string))
	}
	if d.HasChange("ttl") {
		params.Ttl = aws.Int64(d.Get("ttl").(int64))
	}
	if d.HasChange("type") {
		params.Type = aws.String(d.Get("type").(string))
	}

	_, err := conn.UpdateApiCache(params)
	if err != nil {
		return err
	}

	return resourceAwsAppsyncApiCacheRead(d, meta)
}

func resourceAwsAppsyncApiCacheDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appsyncconn

	apiID := d.Id()

	input := &appsync.DeleteApiCacheInput{
		ApiId: aws.String(apiID),
	}
	_, err := conn.DeleteApiCache(input)
	if err != nil {
		if isAWSErr(err, appsync.ErrCodeNotFoundException, "") {
			return nil
		}
		return err
	}

	return nil
}
