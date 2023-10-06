package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/JitenPalaparthi/atipaday/interfaces"
	"github.com/JitenPalaparthi/atipaday/models"
	"github.com/JitenPalaparthi/dapr-go-http-wrapper/wrapper"
	"github.com/golang/glog"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Tip struct {
	ITip interfaces.ITip
	Dapr *wrapper.Dapr
}

func (cp *Tip) Create(ctx context.Context) func(*gin.Context) {
	return func(c *gin.Context) {
		tip := new(models.Tip)
		data := make(map[string]any)
		err := c.ShouldBindBodyWith(&data, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": err.Error(), "message": "bad request"})
			return
		}
		err = c.ShouldBindBodyWith(tip, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": err.Error(), "message": "bad request"})
			return
		}

		tip.Tags = tip.ToString() // adding tags

		if tip.Status == "" {
			tip.Status = "inactive"
		}
		tip.LastModified = time.Now().Unix()
		if con, err := cp.ITip.Create(ctx, tip); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": err.Error(), "message": "Error in creating contact info"})
			return
		} else {
			err := cp.Dapr.Publish("tipspubsub", "tips-created", con.ToBytes())
			if err != nil {
				glog.Errorln(err)
			}
			c.JSON(http.StatusCreated, con)
			return
		}
	}
}

func (cp *Tip) GetBy(ctx context.Context) func(*gin.Context) {
	return func(c *gin.Context) {
		id, ok := c.Params.Get("id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": "invalid id parameter", "message": "invalid id parameter"})
			return
		}

		tip, err := cp.ITip.GetBy(ctx, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": err.Error(), "message": "bad request"})
			return
		}
		c.JSON(http.StatusOK, tip)
	}
}

func (cp *Tip) GetAllByOffset(ctx context.Context) func(*gin.Context) {
	return func(c *gin.Context) {

		limit, ok := c.Params.Get("limit")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": "invalid id parameter", "message": "invalid id parameter"})
			return
		}
		_limit, err := strconv.Atoi(limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": err.Error(), "message": "bad request"})
			return
		}
		offset, ok := c.Params.Get("offset")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": "invalid id parameter", "message": "invalid id parameter"})
			return
		}
		_offset, err := strconv.Atoi(offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": err.Error(), "message": "bad request"})
			return
		}
		tips, err := cp.ITip.GetAllByOffSet(ctx, _offset, _limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": err.Error(), "message": "bad request"})
			return
		}
		c.JSON(http.StatusOK, tips)
	}
}

func (cp *Tip) Search(ctx context.Context) func(*gin.Context) {
	return func(c *gin.Context) {

		limit, ok := c.Params.Get("limit")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": "invalid id parameter", "message": "invalid id parameter"})
			return
		}
		_limit, err := strconv.Atoi(limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": err.Error(), "message": "bad request"})
			return
		}
		offset, ok := c.Params.Get("offset")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": "invalid id parameter", "message": "invalid id parameter"})
			return
		}
		_offset, err := strconv.Atoi(offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": err.Error(), "message": "bad request"})
			return
		}

		search, ok := c.Params.Get("search")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": "invalid id parameter", "message": "invalid id parameter"})
			return
		}

		tips, err := cp.ITip.Search(ctx, _offset, _limit, search)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": err.Error(), "message": "bad request"})
			return
		}
		c.JSON(http.StatusOK, tips)
	}
}

func (cp *Tip) DeleteBy(ctx context.Context) func(*gin.Context) {
	return func(c *gin.Context) {
		id, ok := c.Params.Get("id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": "invalid id parameter", "message": "invalid id parameter"})
			return
		}

		noOfRecords, err := cp.ITip.DeleteBy(ctx, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": err.Error(), "message": "bad request"})
			return
		}
		c.JSON(http.StatusAccepted, noOfRecords)
	}
}

func (cp *Tip) UpdateBy(ctx context.Context) func(*gin.Context) {
	return func(c *gin.Context) {
		id, ok := c.Params.Get("id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": "invalid id parameter", "message": "invalid id parameter"})
			return
		}

		data := make(map[string]any)
		err := c.Bind(&data)
		if err != nil {
			if len(data) == 0 || data == nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": "no filed to edit and update", "message": "bad request"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": err.Error(), "message": "bad request"})
			return
		}

		Tip, err := cp.ITip.UpdateBy(ctx, id, data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": 400, "innerError": err.Error(), "message": "bad request"})
			return
		}
		c.JSON(http.StatusOK, Tip)
	}
}
