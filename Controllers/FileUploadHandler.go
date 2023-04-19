package Controllers

import (
	"fmt"
	"io"
	response "main/Response"
	"net/http"
	"os"
	"path/filepath"
	"time"
	db "main/Database"
	models "main/Models"
	constant "main/Utils"
)

func FileUpload(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.ShowResponse(
			"Failure",
			400,
			"File upload failed",
			err.Error(),
			w,
		)
		return
	}

	//file size check
	r.Body = http.MaxBytesReader(w, r.Body, constant.MAX_UPLOADED_SIZE)
	if err := r.ParseMultipartForm(constant.MAX_UPLOADED_SIZE); err != nil {
		response.ShowResponse(
			"Failure",
			400,
			"File size must not be greater than 40mb",
			err.Error(),
			w,
		)
		return
	}

	defer file.Close()
	
	fmt.Println("file: ", fileHeader)
	
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		response.ShowResponse(
			"Failure",
			400,
			"Failed reading file",
			err.Error(),
			w,
		)
		return
	}

	filetype := http.DetectContentType(buff)
	fmt.Println("filetype: ", filetype)
	if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/gif" && filetype != "application/pdf" && filetype != "application/msword" && filetype != "application/zip" && filetype != "application/octet-stream" && filetype != "video/mp4" && filetype != "audio/mpeg" && filetype != "audio/wave" && filetype != "application/ogg" {
		response.ShowResponse(
			"Failure",
			400,
			"File Format not supported",
			"",
			w,
		)
		return
	}
	
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		response.ShowResponse(
			"Failure",
			400,
			"Failed to seek file",
			"",
			w,
		)
		return
	}

	userId := r.FormValue("userId")
	if userId == ""{
		response.ShowResponse(
			"Failure",
			400,
			"UserId is required",
			"",
			w,
		)
		return
	}

	roomId := r.FormValue("roomId")
	if roomId == ""{
		response.ShowResponse(
			"Failure",
			400,
			"Room Id is required",
			"",
			w,
		)
		return
	}

	var userExists bool
	db.DB.Raw("SELECT EXISTS(SELECT * FROM users where user_id = ? and is_active = true and is_deleted = false)" , userId).Scan(&userExists)
	if !userExists{
		response.ShowResponse(
			"Failure",
			400,
			"User is either not present or is logged out",
			"",
			w,
		)
		return
	}

	var roomExists bool
	db.DB.Raw("SELECT EXISTS(SELECT * FROM rooms where room_id = ? and is_deleted = false)" , roomId).Scan(&roomExists)
	if !roomExists{
		response.ShowResponse(
			"Failure",
			400,
			"Room is either not present or is deleted",
			"",
			w,
		)
		return
	}

	
	var participantExists bool
	db.DB.Raw("SELECT EXISTS(SELECT * FROM participants WHERE room_id = ? and user_id = ? and has_left = false)" , roomId , userId).Scan(&participantExists)
	if !participantExists{
		response.ShowResponse(
			"Failure",
			400,
			"User is not a participant of this room",
			"",
			w,
		)
		return
	}


	err = os.MkdirAll("./File/uploads", os.ModePerm)
	if err != nil {
		response.ShowResponse(
			"Failure",
			500,
			"Failed making directory",
			err.Error(),
			w,
		)
		return
	}

	dst, err := os.Create(fmt.Sprintf("./File/uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))

	if err != nil {
		response.ShowResponse(
			"Failure",
			500,
			"Failed assigning filename",
			err.Error(),
			w,
		)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		response.ShowResponse(
			"Failure",
			500,
			"Failed copying file name",
			err.Error(),
			w,
		)
		return
	}

	//TODO file size limitation
	
	var message models.Message
	
	message.Room_id = roomId
	message.User_id = userId
	message.Text = fileHeader.Filename
	messageType := MessageTypeSet(filetype)
	message.Message_type = messageType
	db.DB.Create(&message)

	response.ShowResponse(
		"Success",
		200,
		"File uploaded successfully",
		message,
		w,
	)
}

func MessageTypeSet(filetype string)string{
	var messageType string

	if filetype == "application/pdf" || filetype == "application/msword" || filetype == "application/zip" || filetype == "application/octet-stream"{
		messageType = "Document"
	}

	if filetype == "image/jpeg" || filetype == "image/png" || filetype == "image/gif"{
		messageType = "Image"
	}

	if filetype == "video/mp4"{
		messageType = "Video"
	}

	if filetype == "audio/mpeg" || filetype == "audio/wave" || filetype == "application/ogg"{
		messageType = "Audio"
	}


	return messageType
}
